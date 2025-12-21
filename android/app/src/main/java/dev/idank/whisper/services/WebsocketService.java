package dev.idank.whisper.services;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.Intent;
import android.os.Build;
import android.os.IBinder;

import androidx.annotation.Nullable;
import androidx.core.app.NotificationCompat;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import dev.idank.whisper.R;
import dev.idank.whisper.clients.WebsocketClient;

public class WebsocketService extends Service {

    private final ExecutorService executorService = Executors.newSingleThreadExecutor();
    private WebsocketClient wsClient;

    @Override
    public void onCreate() {
        super.onCreate();
        startForeground(1, createNotification());
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {

        if (intent != null && wsClient == null) {
            String apiToken = intent.getStringExtra("apiToken");
            String websocketURL = intent.getStringExtra("websocketURL");
            String deviceID = intent.getStringExtra("deviceID");

            WebsocketClient.setInstance(this, apiToken, websocketURL, deviceID);
            wsClient = WebsocketClient.getInstance();

            executorService.submit(wsClient::connect);
        }

        return START_STICKY;
    }

    @Override
    public void onDestroy() {
        super.onDestroy();
        if (wsClient != null) {
            executorService.submit(wsClient::shutdown);
            executorService.shutdown();
        }
    }

    @Nullable
    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }

    private Notification createNotification() {
        String channelId = "whisper_service_channel";

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel(
                    channelId,
                    "Service Channel",
                    NotificationManager.IMPORTANCE_LOW
            );
            getSystemService(NotificationManager.class).createNotificationChannel(channel);
        }

        return new NotificationCompat.Builder(this, channelId)
                .setContentTitle("Service Running")
                .setContentText("Background task active")
                .setSmallIcon(R.drawable.ic_launcher_foreground)
                .build();
    }
}
