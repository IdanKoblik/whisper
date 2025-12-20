package dev.idank.application;

import static dev.idank.application.MainActivity.APP_PREFS;
import static dev.idank.application.MainActivity.WS_PREF;

import android.Manifest;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.os.Build;
import android.os.Bundle;
import android.util.Log;
import android.widget.EditText;
import android.widget.TextView;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.content.ContextCompat;

import com.google.android.material.materialswitch.MaterialSwitch;

import dev.idank.application.listeners.WsTextListener;
import dev.idank.application.services.WebsocketService;

public class HomeActivity extends AppCompatActivity {

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            if (checkSelfPermission(Manifest.permission.POST_NOTIFICATIONS) != PackageManager.PERMISSION_GRANTED) {

                requestPermissions(
                        new String[]{Manifest.permission.POST_NOTIFICATIONS},
                        101
                );
            }

            if (ContextCompat.checkSelfPermission(this, Manifest.permission.SEND_SMS)
                    != PackageManager.PERMISSION_GRANTED) {

                requestPermissions(
                        new String[]{Manifest.permission.SEND_SMS},
                        102
                );
            }
        }

        setContentView(R.layout.activity_home);
        final String apiToken = getIntent().getStringExtra("apiToken");
        if (apiToken == null || apiToken.isEmpty()) {
            Log.e("API_TOKEN", "Api token not found");
            return;
        }

        final String deviceId = getIntent().getStringExtra("deviceId");
        if (deviceId == null || deviceId.isEmpty()) {
            Log.e("DEVICE_ID", "Device id not found");
            return;
        }

        TextView txtApiToken = findViewById(R.id.txtApiToken);
        TextView txtDeviceId = findViewById(R.id.txtDeviceId);

        txtApiToken.setText(mask(apiToken));
        txtDeviceId.setText(deviceId);

        EditText edtWsUrl = findViewById(R.id.edtWsUrl);
        SharedPreferences prefs = getSharedPreferences(APP_PREFS, MODE_PRIVATE);
        String savedUrl = getIntent().getStringExtra(WS_PREF);
        edtWsUrl.setText(savedUrl);

        edtWsUrl.addTextChangedListener(new WsTextListener(prefs));
        TextView txtServiceStatus = findViewById(R.id.txtServiceStatus);

        MaterialSwitch switchService = findViewById(R.id.switchService);
        switchService.setOnCheckedChangeListener((buttonView, isChecked) -> {
            if (isChecked) {
                startServiceTask(txtServiceStatus);
            } else {
                stopServiceTask(txtServiceStatus);
            }
        });

        switchService.setChecked(true);
    }

    private void startServiceTask(TextView txtServiceStatus) {
        Intent serviceIntent = new Intent(this, WebsocketService.class);
        EditText edtWsUrl = findViewById(R.id.edtWsUrl);

        serviceIntent.putExtra("apiToken", getIntent().getStringExtra("apiToken"));
        serviceIntent.putExtra("websocketURL", edtWsUrl.getText().toString());
        serviceIntent.putExtra("deviceID", getIntent().getStringExtra("deviceId"));
        ContextCompat.startForegroundService(this, serviceIntent);

        txtServiceStatus.setText("Service running");
    }

    private void stopServiceTask(TextView txtServiceStatus) {
        Intent serviceIntent = new Intent(this, WebsocketService.class);
        stopService(serviceIntent);
        txtServiceStatus.setText("Service stopped");
    }

    private String mask(String str) {
        return new String(new char[str.length()]).replace("\0", "*");
    }
}
