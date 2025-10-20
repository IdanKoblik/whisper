package me.idank.whisper.managers

import me.idank.whisper.WsListener
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import okio.ByteString
import org.json.JSONObject
import java.util.concurrent.TimeUnit

class WsManager(private var url: String) {

    private var headers: Map<String, String> = emptyMap()

    private val client = OkHttpClient.Builder()
        .readTimeout(0, TimeUnit.MILLISECONDS)
        .build()

    private var ws: WebSocket? = null

    fun setHeaders(newHeaders: Map<String, String>) {
        headers = newHeaders
    }

    fun connect(listener: WsListener) {
        disconnect()
        val requestBuilder = Request.Builder().url(url)
        headers.forEach { (key, value) ->
            requestBuilder.addHeader(key, value)
        }

        val request = requestBuilder.build()

        ws = client.newWebSocket(request, listener)
    }

    fun sendJson(json: JSONObject) {
        ws?.send(json.toString())
    }

    fun disconnect() {
        ws?.close(1000, "Closed by user")
        ws = null
    }
}
