package dev.idank.application;

import static androidx.core.app.ActivityCompat.requestPermissions;

import android.Manifest;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.os.Build;
import android.os.Bundle;

import androidx.activity.EdgeToEdge;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.graphics.Insets;
import androidx.core.view.ViewCompat;
import androidx.core.view.WindowInsetsCompat;

public class MainActivity extends AppCompatActivity {

    public static final String
            APP_PREFS = "WHISPER_PREFS",
            WS_PREF = "websocketURL";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        EdgeToEdge.enable(this);
        setContentView(R.layout.activity_home);

        Intent intent = new Intent(MainActivity.this, HomeActivity.class);
        intent.putExtra("apiToken", "");
        intent.putExtra("deviceId", "");
        intent.putExtra("websocketURL", "");

        startActivity(intent);
    }
}