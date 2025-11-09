package me.idank.whisper

import android.content.Context
import androidx.compose.runtime.*
import androidx.compose.ui.platform.LocalContext
import java.util.*
import androidx.core.content.edit
import java.time.Instant

@Composable
fun AppRoot() {
    val context = LocalContext.current

    var apiToken by remember {
        mutableStateOf(context.getSavedApiToken() ?: "")
    }
    var isLoggedIn by remember {
        mutableStateOf(apiToken.isNotBlank())
    }

    var deviceId by remember { mutableStateOf(
        context.getSavedDeviceID() ?: (UUID.randomUUID().toString() + Instant.now())
    ) }

    if (!isLoggedIn) {
        LoginScreen { key ->
            if (key.isNotBlank()) {
                apiToken = key
                context.saveLogin(key, deviceId)
            }
        }
    } else {
        HomeScreen(apiToken, deviceId) {
            context.clearLogin()
        }
    }
}

fun Context.saveLogin(apiToken: String, deviceID: String) {
    val prefs = getSharedPreferences("whisper_prefs", Context.MODE_PRIVATE)
    prefs.edit {
        putString("api_token", apiToken)
            .putString("device_id", deviceID)
            .putBoolean("is_logged_in", true)
    }
}

fun Context.clearLogin() {
    val prefs = getSharedPreferences("whisper_prefs", Context.MODE_PRIVATE)
    prefs.edit {
        remove("api_token")
            .remove("device_id")
            .putBoolean("is_logged_in", false)
    }
}

fun Context.getSavedApiToken(): String? {
    val prefs = getSharedPreferences("whisper_prefs", Context.MODE_PRIVATE)
    return if (prefs.getBoolean("is_logged_in", false)) {
        prefs.getString("api_token", null)
    } else null
}

fun Context.getSavedDeviceID(): String? {
    val prefs = getSharedPreferences("whisper_prefs", Context.MODE_PRIVATE)
    return if (prefs.getBoolean("is_logged_in", false)) {
        prefs.getString("device_id", null)
    } else null
}