package me.idank.whisper

import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import okio.ByteString

class WsListener(
    val onOpen: (() -> Unit)? = null,
    val onMessage: ((String) -> Unit)? = null,
    val onClose: (() -> Unit)? = null,
    val onError: ((Throwable) -> Unit)? = null
) : WebSocketListener() {

    override fun onOpen(webSocket: WebSocket, response: Response) {
        onOpen?.invoke()
    }

    override fun onMessage(webSocket: WebSocket, text: String) {
        onMessage?.invoke(text)
    }

    override fun onMessage(webSocket: WebSocket, bytes: ByteString) {
        onMessage?.invoke(bytes.utf8())
    }

    override fun onClosing(webSocket: WebSocket, code: Int, reason: String) {
        webSocket.close(1000, null)
        onClose?.invoke()
    }

    override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
        t.printStackTrace()
        onError?.invoke(t)
    }
}