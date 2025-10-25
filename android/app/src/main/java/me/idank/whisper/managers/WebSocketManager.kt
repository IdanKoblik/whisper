package me.idank.whisper.managers

import android.content.Context
import android.util.Log
import me.idank.whisper.listeners.ReconnectingWebSocketListener
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import java.util.concurrent.TimeUnit

class WebSocketManager private constructor(
    var serverUrl: String,
    var apiToken: String,
    var deviceID: String
) {
    private var client: OkHttpClient? = null
    private var webSocket: WebSocket? = null
    private var listener: WebSocketListener? = null

    companion object {
        @Volatile
        private var INSTANCE: WebSocketManager? = null

        fun getInstance(serverUrl: String, apiToken: String, deviceID: String): WebSocketManager {
            return INSTANCE ?: synchronized(this) {
                INSTANCE ?: WebSocketManager(serverUrl, apiToken, deviceID).also { INSTANCE = it }
            }
        }

        fun getInstance(): WebSocketManager {
            return INSTANCE
                ?: throw IllegalStateException("WebSocketManager not initialized. Call getInstance(serverUrl, apiToken) first.")
        }
    }

    fun connectWithContext(context: Context) {
        if (client == null) {
            client = OkHttpClient.Builder()
                .readTimeout(0, TimeUnit.MILLISECONDS)
                .build()
        }

        listener = ReconnectingWebSocketListener(context, this)

        val request = Request.Builder()
            .url(serverUrl)
            .addHeader("X-Api-Token", apiToken)
            .build()

        webSocket = client?.newWebSocket(request, listener!!)
        Log.d("WebSocketManager", "Connecting to $serverUrl...")
    }

    fun reconnectWithContext(context: Context) {
        disconnect()
        connectWithContext(context)
    }

    fun disconnect() {
        webSocket?.close(1000, "Client closed")
        webSocket = null
        listener = null
        Log.d("WebSocketManager", "Disconnected.")
    }

    fun shutdown() {
        disconnect()
        client?.dispatcher?.executorService?.shutdown()
        client = null
        Log.d("WebSocketManager", "Client shutdown.")
    }
}
