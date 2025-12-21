package dev.idank.whisper.listeners;

import static dev.idank.whisper.MainActivity.WS_PREF;

import android.content.SharedPreferences;
import android.text.Editable;
import android.text.TextWatcher;

public class WsTextListener implements TextWatcher {

    private final SharedPreferences prefs;

    public WsTextListener(SharedPreferences prefs) {
        this.prefs = prefs;
    }

    @Override
    public void afterTextChanged(Editable s) {
        this.prefs.edit().putString(WS_PREF, s.toString()).apply();
    }

    @Override
    public void beforeTextChanged(CharSequence s, int start, int count, int after) {
        // Nothing
    }

    @Override
    public void onTextChanged(CharSequence s, int start, int before, int count) {
        // Nothing
    }
}
