package me.idank.whisper

import android.content.Intent
import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Close
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
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
    var wsUrl by remember { mutableStateOf("") }
    var serviceActive by remember { mutableStateOf(false) }
    var showKey by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Whisper") },
                actions = {
                    TextButton(onClick = onLogout) {
                        Text("Logout")
                    }
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
            // API Token Card
            Card {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("API token", style = MaterialTheme.typography.labelLarge)
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(
                            text = if (showKey) apiToken else "•".repeat(apiToken.length.coerceAtLeast(6)),
                            style = MaterialTheme.typography.bodyMedium
                        )
                        Spacer(Modifier.width(8.dp))
                        IconButton(onClick = { showKey = !showKey }) {
                            Icon(
                                imageVector = if (showKey) Icons.Default.Close else Icons.Default.Lock,
                                contentDescription = "Toggle token visibility"
                            )
                        }
                    }
                }
            }

            // Device ID Card
            Card {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Device ID", style = MaterialTheme.typography.labelLarge)
                    Text(deviceId, style = MaterialTheme.typography.bodyMedium)
                }
            }

            // WebSocket URL Card
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

            // Service Toggle Card
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
                            val intent = Intent(context, WebSocketService::class.java)
                            if (active) {
                                WebSocketManager.getInstance(wsUrl, apiToken, deviceId)
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