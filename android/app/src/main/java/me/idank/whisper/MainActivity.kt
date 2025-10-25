package me.idank.whisper

import android.annotation.SuppressLint
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.os.Bundle
import android.telephony.SmsManager
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.core.app.NotificationCompat
import androidx.core.app.NotificationManagerCompat
import me.idank.whisper.ui.theme.AppTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        createNotificationChannel()
        setContent {
            AppTheme(darkTheme = true, dynamicColor = false) {
                AppRoot()
            }
        }
    }
}

fun Context.createNotificationChannel() {
    val channel = NotificationChannel(
        "whisper_channel",
        "Whisper Notifications",
        NotificationManager.IMPORTANCE_DEFAULT
    ).apply { description = "Notifications from Whisper app" }

    val manager = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
    manager.createNotificationChannel(channel)
}

@SuppressLint("MissingPermission")
fun Context.showNotification(title: String, message: String) {
    val intent = Intent(this, MainActivity::class.java).apply {
        flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
    }

    val pendingIntent = PendingIntent.getActivity(
        this,
        0,
        intent,
        PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
    )

    val builder = NotificationCompat.Builder(this, "whisper_channel")
        .setSmallIcon(R.drawable.ic_launcher_background)
        .setContentTitle(title)
        .setContentText(message)
        .setPriority(NotificationCompat.PRIORITY_DEFAULT)
        .setContentIntent(pendingIntent)
        .setAutoCancel(true)

    NotificationManagerCompat.from(this).notify((System.currentTimeMillis() % Int.MAX_VALUE).toInt(), builder.build())
}

fun Context.sendSms(phoneNumber: String, message: String) {
    try {
        SmsManager.getDefault().sendTextMessage(phoneNumber, null, message, null, null)
        showNotification("SMS sent", "Message sent to $phoneNumber")
    } catch (e: Exception) {
        showNotification("SMS failed", "Failed to send to $phoneNumber: ${e.message}")
    }
}
