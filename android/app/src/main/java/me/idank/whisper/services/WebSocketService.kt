package me.idank.whisper.services

import android.R
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.Service
import android.content.Intent
import android.os.IBinder
import androidx.core.app.NotificationCompat
import me.idank.whisper.managers.WebSocketManager

class WebSocketService : Service() {

    private lateinit var manager: WebSocketManager

    override fun onCreate() {
        super.onCreate()
        startForegroundServiceNotification()

        manager = WebSocketManager.getInstance()
        manager.connectWithContext(this)
    }

    private fun startForegroundServiceNotification() {
        val channelId = "websocket_channel"
        val channel = NotificationChannel(
            channelId,
            "WebSocket Service",
            NotificationManager.IMPORTANCE_LOW
        )
        (getSystemService(NOTIFICATION_SERVICE) as NotificationManager)
            .createNotificationChannel(channel)

        val notification = NotificationCompat.Builder(this, channelId)
            .setContentTitle("Whisper WebSocket Connected")
            .setSmallIcon(R.drawable.stat_sys_download_done)
            .build()

        startForeground(1, notification)
    }

    override fun onDestroy() {
        super.onDestroy()
        manager.disconnect()
    }

    override fun onBind(intent: Intent?): IBinder? = null
}
