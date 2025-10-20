package me.idank.whisper

import android.Manifest
import android.annotation.SuppressLint
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Build
import android.os.Bundle
import android.telephony.SmsManager
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.core.app.ActivityCompat
import androidx.core.app.NotificationCompat
import androidx.core.app.NotificationManagerCompat
import androidx.core.content.ContextCompat
import me.idank.whisper.managers.WsManager
import me.idank.whisper.ui.theme.AppTheme
import org.json.JSONObject
import java.util.*

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

    with(NotificationManagerCompat.from(this)) {
        notify((System.currentTimeMillis() % Int.MAX_VALUE).toInt(), builder.build())
    }
}

@Composable
fun AppRoot() {
    var signatureKey by remember { mutableStateOf("") }
    var isLoggedIn by remember { mutableStateOf(false) }
    var deviceId by remember { mutableStateOf(UUID.randomUUID().toString()) }

    if (!isLoggedIn) {
        LoginScreen { key ->
            if (key.isNotBlank()) {
                signatureKey = key
                isLoggedIn = true
            }
        }
    } else {
        HomeScreen(signatureKey, deviceId) {
            signatureKey = ""
            isLoggedIn = false
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LoginScreen(onSubmit: (String) -> Unit) {
    var keyInput by remember { mutableStateOf("") }

    Scaffold(
        topBar = { TopAppBar(title = { Text("Whisper") }) }
    ) { innerPadding ->
        Column(
            modifier = Modifier
                .padding(innerPadding)
                .padding(24.dp)
                .fillMaxSize(),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text("Sign in", style = MaterialTheme.typography.headlineMedium)
            Spacer(Modifier.size(24.dp))
            OutlinedTextField(
                value = keyInput,
                onValueChange = { keyInput = it },
                label = { Text("Signature key") },
                singleLine = true,
                visualTransformation = PasswordVisualTransformation(),
                keyboardOptions = KeyboardOptions(autoCorrect = false),
                modifier = Modifier.fillMaxWidth()
            )
            Spacer(Modifier.size(16.dp))
            Button(
                onClick = { onSubmit(keyInput) },
                modifier = Modifier.fillMaxWidth()
            ) {
                Text("Continue")
            }
        }
    }
}


fun Context.sendSms(phoneNumber: String, message: String) {
    try {
        val smsManager = SmsManager.getDefault()
        smsManager.sendTextMessage(phoneNumber, null, message, null, null)
        showNotification("SMS sent", "Message sent to $phoneNumber")
    } catch (e: Exception) {
        showNotification("SMS failed", "Failed to send to $phoneNumber: ${e.message}")
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(signatureKey: String, deviceId: String, onLogout: () -> Unit) {
    val context = LocalContext.current
    var wsUrl by remember { mutableStateOf("") }
    var serviceActive by remember { mutableStateOf(false) }
    var connectionState by remember { mutableStateOf("Disconnected") }

    val wsManager = remember(wsUrl) { WsManager(wsUrl) }
    wsManager.setHeaders(mapOf("X-Api-Token" to signatureKey))

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Whisper") },
                actions = {
                    TextButton(onClick = onLogout) { Text("Logout") }
                }
            )
        }
    ) { innerPadding ->
        Column(
            modifier = Modifier
                .padding(innerPadding)
                .padding(24.dp)
                .fillMaxSize(),
            verticalArrangement = Arrangement.spacedBy(24.dp)
        ) {
            Card {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Signature key", style = MaterialTheme.typography.labelLarge)
                    Text("•".repeat(signatureKey.length.coerceAtLeast(6)))
                }
            }

            Card {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Device ID", style = MaterialTheme.typography.labelLarge)
                    Text(deviceId)
                }
            }

            Card {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("WebSocket URL", style = MaterialTheme.typography.labelLarge)
                    OutlinedTextField(
                        value = wsUrl,
                        onValueChange = { wsUrl = it },
                        placeholder = { Text("Enter WebSocket URL") },
                        singleLine = true,
                        modifier = Modifier.fillMaxWidth()
                    )
                }
            }

            Card {
                Row(
                    modifier = Modifier
                        .padding(16.dp)
                        .fillMaxWidth(),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Column {
                        Text("Activate service", style = MaterialTheme.typography.bodyLarge)
                        Text(
                            text = "Status: $connectionState",
                            style = MaterialTheme.typography.bodyMedium
                        )
                    }

                    Switch(
                        checked = serviceActive,
                        onCheckedChange = { active ->
                            serviceActive = active
                            try {
                                if (active) {
                                    if (wsUrl.isBlank()) {
                                        context.showNotification(
                                            "Whisper",
                                            "Please enter a valid WebSocket URL."
                                        )
                                        serviceActive = false
                                        return@Switch
                                    }

                                    connectionState = "Connecting..."
                                    wsManager.connect(
                                        WsListener(
                                            onOpen = {
                                                connectionState = "Connected"
                                                val json = JSONObject().put("device_id", deviceId)
                                                wsManager.sendJson(json)
                                                context.showNotification(
                                                    "Whisper",
                                                    "Service activated successfully!"
                                                )
                                            },
                                            onMessage = { message ->
                                                context.showNotification("Message received", message)
                                                try {
                                                    val json = JSONObject(message)
                                                    val subscribers = mutableListOf<String>()
                                                    val array = json.optJSONArray("subscribers")
                                                    if (array != null) {
                                                        for (i in 0 until array.length()) {
                                                            val phone = array.optString(i)
                                                            if (phone.isNotBlank()) subscribers.add(phone)
                                                        }
                                                    }

                                                    subscribers.forEach { phone ->
                                                        context.sendSms(phone, "You received a message from Whisper!")
                                                    }

                                                } catch (e: Exception) {
                                                    context.showNotification(
                                                        "Parsing error",
                                                        "Failed to parse subscribers: ${e.message}"
                                                    )
                                                }
                                            },
                                            onClose = {
                                                connectionState = "Disconnected"
                                            },
                                            onError = { err ->
                                                connectionState = "Error: ${err.message}"
                                                context.showNotification(
                                                    "Whisper",
                                                    "WebSocket error: ${err.message}"
                                                )
                                                serviceActive = false
                                            }
                                        )
                                    )
                                } else {
                                    val json = JSONObject().put("device_id", deviceId)
                                    wsManager.sendJson(json)
                                    wsManager.disconnect()
                                    connectionState = "Disconnected"
                                    context.showNotification("Whisper", "Service deactivated!")
                                }
                            } catch (e: Exception) {
                                connectionState = "Error: ${e.message}"
                                context.showNotification("Whisper", "Error: ${e.message}")
                            }
                        }
                    )
                }
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun LoginPreview() {
    AppTheme(darkTheme = true, dynamicColor = false) {
        LoginScreen(onSubmit = {})
    }
}

@Preview(showBackground = true)
@Composable
fun HomePreview() {
    AppTheme(darkTheme = true, dynamicColor = false) {
        HomeScreen(signatureKey = "sk_live_1234567890", deviceId = "device_123", onLogout = {})
    }
}
