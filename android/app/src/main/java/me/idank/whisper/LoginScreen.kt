package me.idank.whisper

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Close
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LoginScreen(onSubmit: (String) -> Unit) {
    var keyInput by remember { mutableStateOf("") }
    var passwordVisible by remember { mutableStateOf(false) }

    Scaffold(topBar = { TopAppBar(title = { Text("Whisper") }) }) { innerPadding ->
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
                label = { Text("API token") },
                singleLine = true,
                visualTransformation = if (passwordVisible) VisualTransformation.None else PasswordVisualTransformation(),
                trailingIcon = { },
                keyboardOptions = KeyboardOptions(autoCorrect = false),
                modifier = Modifier.fillMaxWidth()
            )
            Spacer(Modifier.size(16.dp))
            Button(onClick = { onSubmit(keyInput) }, modifier = Modifier.fillMaxWidth()) {
                Text("Continue")
            }
        }
    }
}

