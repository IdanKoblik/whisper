package dev.idank.application;

import static dev.idank.application.MainActivity.APP_PREFS;
import static dev.idank.application.MainActivity.WS_PREF;
import static dev.idank.application.MainActivity.API_TOKEN_PREF;

import android.Manifest;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.os.Build;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.EditText;
import android.widget.TextView;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.content.ContextCompat;

import com.google.android.material.button.MaterialButton;
import com.google.android.material.materialswitch.MaterialSwitch;

import dev.idank.application.listeners.WsTextListener;
import dev.idank.application.services.WebsocketService;

public class HomeActivity extends AppCompatActivity {

    private String apiToken;
    private boolean isTokenVisible = false;
    private TextView txtApiToken;
    private MaterialButton btnToggleToken;

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            if (checkSelfPermission(Manifest.permission.POST_NOTIFICATIONS) != PackageManager.PERMISSION_GRANTED) {
                requestPermissions(new String[]{Manifest.permission.POST_NOTIFICATIONS}, 101);
            }

            if (ContextCompat.checkSelfPermission(this, Manifest.permission.SEND_SMS) != PackageManager.PERMISSION_GRANTED) {
                requestPermissions(new String[]{Manifest.permission.SEND_SMS}, 102);
            }

            PackageManager pm = getPackageManager();
            if (!pm.hasSystemFeature(PackageManager.FEATURE_TELEPHONY_MESSAGING)) {
                Log.e("HomeActivity", "Device does not support SMS");
                return;
            }
        }

        setContentView(R.layout.activity_home);

        apiToken = getIntent().getStringExtra("apiToken");
        if (apiToken == null || apiToken.isEmpty()) {
            Log.e("HomeActivity", "Api token not found");
            return;
        }

        String deviceId = getIntent().getStringExtra("deviceId");
        if (deviceId == null || deviceId.isEmpty()) {
            Log.e("HomeActivity", "Device id not found");
            return;
        }

        txtApiToken = findViewById(R.id.txtApiToken);
        TextView txtDeviceId = findViewById(R.id.txtDeviceId);
        txtApiToken.setText(mask(apiToken));
        txtDeviceId.setText(deviceId);

        btnToggleToken = findViewById(R.id.btnToggleToken);
        btnToggleToken.setOnClickListener(this::onToggleTokenClick);

        EditText edtWsUrl = findViewById(R.id.edtWsUrl);
        SharedPreferences prefs = getSharedPreferences(APP_PREFS, MODE_PRIVATE);
        String savedUrl = getIntent().getStringExtra(WS_PREF);
        if (savedUrl == null || savedUrl.isEmpty()) {
            savedUrl = prefs.getString(WS_PREF, "");
        }
        edtWsUrl.setText(savedUrl);

        edtWsUrl.addTextChangedListener(new WsTextListener(prefs));

        MaterialButton btnLogout = findViewById(R.id.btnLogout);
        btnLogout.setOnClickListener(this::onLogoutClick);

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

    private void onToggleTokenClick(View view) {
        isTokenVisible = !isTokenVisible;
        if (isTokenVisible) {
            txtApiToken.setText(apiToken);
            btnToggleToken.setText("Hide");
        } else {
            txtApiToken.setText(mask(apiToken));
            btnToggleToken.setText("Show");
        }
    }

    private void onLogoutClick(View view) {
        SharedPreferences prefs = getSharedPreferences(APP_PREFS, MODE_PRIVATE);
        SharedPreferences.Editor editor = prefs.edit();
        editor.remove(API_TOKEN_PREF);
        editor.apply();

        Intent intent = new Intent(HomeActivity.this, LoginActivity.class);
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK | Intent.FLAG_ACTIVITY_CLEAR_TASK);
        startActivity(intent);
        finish();
    }

    private String mask(String str) {
        return new String(new char[str.length()]).replace("\0", "*");
    }
}
