package dev.idank.whisper;

import static dev.idank.whisper.MainActivity.APP_PREFS;
import static dev.idank.whisper.MainActivity.API_TOKEN_PREF;
import static dev.idank.whisper.MainActivity.DEVICE_ID_PREF;

import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;

import androidx.activity.EdgeToEdge;
import androidx.appcompat.app.AppCompatActivity;

import com.google.android.material.button.MaterialButton;
import com.google.android.material.textfield.TextInputEditText;
import com.google.android.material.textfield.TextInputLayout;

import java.util.UUID;

public class LoginActivity extends AppCompatActivity {

    private TextInputEditText edtApiToken;
    private TextInputLayout tilApiToken;
    private MaterialButton btnLogin;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        EdgeToEdge.enable(this);
        setContentView(R.layout.activity_login);

        edtApiToken = findViewById(R.id.edtApiToken);
        tilApiToken = findViewById(R.id.tilApiToken);
        btnLogin = findViewById(R.id.btnLogin);

        btnLogin.setOnClickListener(this::onLoginClick);
    }

    private void onLoginClick(View view) {
        String apiToken = edtApiToken.getText() != null ? edtApiToken.getText().toString().trim() : "";

        if (TextUtils.isEmpty(apiToken)) {
            tilApiToken.setError("API token is required");
            return;
        }

        tilApiToken.setError(null);

        SharedPreferences prefs = getSharedPreferences(APP_PREFS, MODE_PRIVATE);
        SharedPreferences.Editor editor = prefs.edit();
        editor.putString(API_TOKEN_PREF, apiToken);

        String deviceId = prefs.getString(DEVICE_ID_PREF, null);
        if (deviceId == null || deviceId.isEmpty()) {
            deviceId = UUID.randomUUID().toString();
            editor.putString(DEVICE_ID_PREF, deviceId);
        }

        editor.apply();

        Intent intent = new Intent(LoginActivity.this, HomeActivity.class);
        intent.putExtra("apiToken", apiToken);
        intent.putExtra("deviceId", deviceId);
        startActivity(intent);
        finish();
    }
}

