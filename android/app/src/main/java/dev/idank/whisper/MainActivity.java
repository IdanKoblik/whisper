package dev.idank.whisper;

import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;

import androidx.activity.EdgeToEdge;
import androidx.appcompat.app.AppCompatActivity;

import java.util.UUID;

public class MainActivity extends AppCompatActivity {

    public static final String
            APP_PREFS = "WHISPER_PREFS",
            WS_PREF = "websocketURL",
            API_TOKEN_PREF = "apiToken",
            DEVICE_ID_PREF = "deviceId";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        EdgeToEdge.enable(this);

        SharedPreferences prefs = getSharedPreferences(APP_PREFS, MODE_PRIVATE);
        String apiToken = prefs.getString(API_TOKEN_PREF, null);
        String deviceId = prefs.getString(DEVICE_ID_PREF, null);

        Intent intent;
        if (apiToken == null || apiToken.isEmpty()) {
            intent = new Intent(MainActivity.this, LoginActivity.class);
        } else {
            if (deviceId == null || deviceId.isEmpty()) {
                deviceId = UUID.randomUUID().toString();
                SharedPreferences.Editor editor = prefs.edit();
                editor.putString(DEVICE_ID_PREF, deviceId);
                editor.apply();
            }
            intent = new Intent(MainActivity.this, HomeActivity.class);
            intent.putExtra("apiToken", apiToken);
            intent.putExtra("deviceId", deviceId);
        }

        startActivity(intent);
        finish();
    }
}