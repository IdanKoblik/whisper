package dev.idank.application.listeners;

import android.app.Activity;
import android.app.PendingIntent;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.pm.PackageManager;
import android.os.Build;
import android.telephony.SmsManager;
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.core.content.ContextCompat;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.util.ArrayList;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

import dev.idank.application.clients.WebsocketClient;
import dev.idank.application.requests.HeartbeatRequest;
import dev.idank.application.requests.MessageRequest;
import okhttp3.Response;
import okhttp3.WebSocket;
import okhttp3.WebSocketListener;

public class WebsocketListener extends WebSocketListener {

    private final Context context;
    private final ObjectMapper mapper = new ObjectMapper();
    private final ScheduledExecutorService scheduler = Executors.newSingleThreadScheduledExecutor();
    private final WebsocketClient client;

    private ScheduledFuture<?> heartbeatFuture;
    private ScheduledFuture<?> reconnectFuture;

    public WebsocketListener(Context context, WebsocketClient client) {
        this.context = context.getApplicationContext();
        this.client = client;
    }

    @Override
    public void onOpen(@NonNull WebSocket webSocket, @NonNull Response response) {
        Log.d("WebsocketListener", "WebSocket connected");

        client.onConnected();
        stopReconnect();

        try {
            String data = mapper.writeValueAsString(new HeartbeatRequest(client.getDeviceID()));
            webSocket.send(data);
        } catch (JsonProcessingException e) {
            Log.e("WebsocketListener", "Connect payload error", e);
        }

        scheduleHeartbeat(webSocket);
    }

    @Override
    public void onMessage(@NonNull WebSocket webSocket, @NonNull String text) {
        try {
            JsonNode root = mapper.readTree(text);
            if (!root.hasNonNull("message") || !root.hasNonNull("targets")) {
                Log.w("WebsocketListener", "Invalid message format");
                return;
            }

            MessageRequest message = mapper.treeToValue(root, MessageRequest.class);

            if (context.checkSelfPermission(android.Manifest.permission.SEND_SMS) != PackageManager.PERMISSION_GRANTED) {
                Log.e("WebsocketListener", "SEND_SMS permission not granted");
                return;
            }

            SmsManager smsManager;
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
                smsManager = context.getSystemService(SmsManager.class);
            } else {
                smsManager = SmsManager.getDefault();
            }

            if (smsManager == null) {
                Log.e("WebsocketListener", "SmsManager not available");
                return;
            }

            for (String target : message.targets()) {
                Log.d("WebsocketListener", "Sending SMS to " + target);
                ArrayList<String> parts = smsManager.divideMessage(message.message());

                smsManager.sendMultipartTextMessage(target, null, parts, null, null);
            }

        } catch (Exception e) {
            Log.e("WebsocketListener", "Invalid message JSON or SMS sending failed", e);
        }
    }

    @Override
    public void onClosed(@NonNull WebSocket webSocket, int code, @NonNull String reason) {
        Log.d("WebsocketListener", "Closed: " + reason);
        client.onDisconnected();
        stopHeartbeat();
        scheduleReconnect();
    }

    @Override
    public void onFailure(@NonNull WebSocket webSocket, @NonNull Throwable t, @Nullable Response response) {
        Log.e("WebsocketListener", "Failure", t);
        client.onDisconnected();
        stopHeartbeat();
        scheduleReconnect();
    }

    private void scheduleHeartbeat(WebSocket webSocket) {
        if (heartbeatFuture != null && !heartbeatFuture.isDone()) return;

        heartbeatFuture = scheduler.scheduleWithFixedDelay(() -> {
            try {
                webSocket.send(mapper.writeValueAsString(new HeartbeatRequest(client.getDeviceID())));
            } catch (JsonProcessingException e) {
                Log.e("WebsocketListener", "Heartbeat send failed", e);
            }
        }, 5, 5, TimeUnit.SECONDS);
    }

    private void stopHeartbeat() {
        if (heartbeatFuture != null) {
            heartbeatFuture.cancel(true);
            heartbeatFuture = null;
        }
    }

    private void scheduleReconnect() {
        if (reconnectFuture != null && !reconnectFuture.isDone()) return;

        reconnectFuture = scheduler.schedule(() -> {
            Log.d("WebsocketListener", "Reconnecting...");
            client.reconnect();
        }, 5, TimeUnit.SECONDS);
    }

    private void stopReconnect() {
        if (reconnectFuture != null) {
            reconnectFuture.cancel(true);
            reconnectFuture = null;
        }
    }
}
