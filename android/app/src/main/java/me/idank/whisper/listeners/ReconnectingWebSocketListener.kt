package me.idank.whisper.listeners

import android.content.Context
import android.util.Log
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import me.idank.whisper.data.MessageRequest
import me.idank.whisper.managers.WebSocketManager
import me.idank.whisper.sendSms
import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import org.json.JSONObject

class ReconnectingWebSocketListener(
    private val context: Context,
    private val manager: WebSocketManager
) : WebSocketListener() {

    private val mapper = jacksonObjectMapper()
    private val scope = CoroutineScope(Dispatchers.IO)
    private var reconnectJob: Job? = null

    override fun onOpen(webSocket: WebSocket, response: Response) {
        Log.d("WebSocket", "Connected")
        Log.d("WebSocket", manager.apiToken)
        reconnectJob?.cancel()

        val json = JSONObject().apply { put("device_id", manager.deviceID) }
        webSocket.send(json.toString())
        Log.d("WebSocket", "Sent JSON body: $json")
    }

    override fun onMessage(webSocket: WebSocket, text: String) {
        Log.d("WebSocket", "Message received: $text")

        try {
            val message = mapper.readValue(text, MessageRequest::class.java)
            Log.d("WebSocket", "Parsed MessageRequest: $message")

            message.subscribers.forEach {
                context.sendSms(it, message.message)
            }

        } catch (e: Exception) {
            Log.w("WebSocket", "Ignored non-MessageRequest: $text")
        }
    }


    override fun onClosed(webSocket: WebSocket, code: Int, reason: String) {
        Log.d("WebSocket", "Closed: $reason")
        scheduleReconnect()
    }

    override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
        Log.e("WebSocket", "Failure: ${t.message}")
        scheduleReconnect()
    }

    private fun scheduleReconnect() {
        if (reconnectJob?.isActive == true)
            return

        reconnectJob = scope.launch {
            delay(5000)
            Log.d("WebSocket", "Reconnecting...")
            manager.reconnectWithContext(context)
        }
    }
}
