package me.idank.whisper

import android.content.Context
import androidx.compose.runtime.*
import androidx.compose.ui.platform.LocalContext
import java.util.*
import androidx.core.content.edit

@Composable
fun AppRoot() {
    val context = LocalContext.current

    var apiToken by remember {
        mutableStateOf(context.getSavedLogin() ?: "68263e6587aded201b805762bfecd18aec608515d08a05f813be7259e9eb35fd")
    }
    var isLoggedIn by remember {
        mutableStateOf(apiToken.isNotBlank())
    }
    var deviceId by remember { mutableStateOf(UUID.randomUUID().toString()) }

    if (!isLoggedIn) {
        LoginScreen { key ->
            if (key.isNotBlank()) {
                apiToken = key
                isLoggedIn = true
                context.saveLogin(key)
            }
        }
    } else {
        HomeScreen(apiToken, deviceId) {
            apiToken = "68263e6587aded201b805762bfecd18aec608515d08a05f813be7259e9eb35fd"
            isLoggedIn = false
            context.clearLogin()
        }
    }
}

fun Context.saveLogin(apiToken: String) {
    val prefs = getSharedPreferences("whisper_prefs", Context.MODE_PRIVATE)
    prefs.edit {
        putString("api_token", apiToken)
            .putBoolean("is_logged_in", true)
    }
}

fun Context.clearLogin() {
    val prefs = getSharedPreferences("whisper_prefs", Context.MODE_PRIVATE)
    prefs.edit {
        remove("api_token")
            .putBoolean("is_logged_in", false)
    }
}

fun Context.getSavedLogin(): String? {
    val prefs = getSharedPreferences("whisper_prefs", Context.MODE_PRIVATE)
    return if (prefs.getBoolean("is_logged_in", false)) {
        prefs.getString("api_token", null)
    } else null
}