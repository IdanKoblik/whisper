package dev.idank.application.listeners;

import android.telephony.SmsManager;
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

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

    private final ObjectMapper mapper = new ObjectMapper();
    private final ScheduledExecutorService scheduler =
            Executors.newSingleThreadScheduledExecutor();

    private final WebsocketClient client;

    private ScheduledFuture<?> heartbeatFuture;
    private ScheduledFuture<?> reconnectFuture;

    public WebsocketListener(WebsocketClient client) {
        this.client = client;
    }

    @Override
    public void onOpen(@NonNull WebSocket webSocket, @NonNull Response response) {
        Log.d(getClass().getName(), "WebSocket connected");

        client.onConnected();
        stopReconnect();

        try {
            String data = mapper.writeValueAsString(
                    new HeartbeatRequest(client.getDeviceID())
            );
            webSocket.send(data);
        } catch (JsonProcessingException e) {
            Log.e(getClass().getName(), "Connect payload error", e);
        }

        scheduleHeartbeat(webSocket);
    }

    @Override
    public void onMessage(@NonNull WebSocket webSocket, @NonNull String text) {
        Log.d(getClass().getName(), "Message: " + text);

        MessageRequest message = null;
        try {
            message = mapper.readValue(text, MessageRequest.class);
        } catch (Exception e) {
            // Pass
        }

        if (message == null)
            return;

        MessageRequest finalMessage = message;
        message.targets().forEach(target -> {
            SmsManager.getDefault().sendTextMessage(target, null, finalMessage.message(), null, null);
        });
    }

    @Override
    public void onClosed(@NonNull WebSocket webSocket, int code, @NonNull String reason) {
        Log.d(getClass().getName(), "Closed: " + reason);

        client.onDisconnected();
        stopHeartbeat();
        scheduleReconnect();
    }

    @Override
    public void onFailure(
            @NonNull WebSocket webSocket,
            @NonNull Throwable t,
            @Nullable Response response
    ) {
        Log.e(getClass().getName(), "Failure", t);

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
                throw new RuntimeException(e);
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
            Log.d(getClass().getName(), "Reconnecting...");
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
