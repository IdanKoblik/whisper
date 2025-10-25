package me.idank.whisper

import android.content.Intent
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import me.idank.whisper.managers.WebSocketManager
import me.idank.whisper.services.WebSocketService

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(apiToken: String, deviceId: String, onLogout: () -> Unit) {
    val context = LocalContext.current
    var wsUrl by remember { mutableStateOf("ws://10.5.0.84:8080/ws") }
    var serviceActive by remember { mutableStateOf(true) }
    var showKey by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Whisper") },
                actions = { TextButton(onClick = onLogout) { Text("Logout") } }
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
                    Text("API token", style = MaterialTheme.typography.labelLarge)
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(if (showKey) apiToken else "•".repeat(apiToken.length.coerceAtLeast(6)))
                        Spacer(Modifier.width(8.dp))
                    }
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
                            text = if (serviceActive) "Service running" else "Service stopped",
                            style = MaterialTheme.typography.bodyMedium
                        )
                    }
                    Switch(
                        checked = serviceActive,
                        onCheckedChange = { active ->
                            serviceActive = active
                            WebSocketManager.getInstance(
                                serverUrl = wsUrl,
                                apiToken = apiToken,
                                deviceID = deviceId
                            )
                            val intent = Intent(context, WebSocketService::class.java)
                            if (active) {
                                context.startForegroundService(intent)
                            } else {
                                context.stopService(intent)
                            }
                        }
                    )
                }
            }
        }
    }
}
