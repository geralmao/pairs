package com.programatta.games.pairs;

import android.os.Build;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.view.WindowCompat;
import androidx.core.view.WindowInsetsCompat;
import androidx.core.view.WindowInsetsControllerCompat;

import java.util.Locale;

import go.Seq;
import com.programatta.games.pairs.corelib.mobile.EbitenView;
import com.programatta.games.pairs.corelib.mobile.Mobile;


public class MainActivity extends AppCompatActivity {
  private static final String TAG = "Pairs"; // Tag for logging

  @Override
  protected void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    this.setContentView(R.layout.activity_main);

    // This ensures your content goes edge-to-edge, behind the system bars
    // It's a key part of modern full-screen layouts
    WindowCompat.setDecorFitsSystemWindows(getWindow(), false);

    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) { // API 30+
        // Usamos WindowInsetsController para Android 11 (API 30) y superior
        hideSystemBarsApi30();
    } else {
        // Para APIs anteriores, seguiríamos usando setSystemUiVisibility()
        hideSystemBarsLegacy();
    }

    // Obtenemos las vistas.
    this.ebitenView = this.getEbitenView();

    //Pasamos a golang la ruta de la aplicación.
    String filesDir = getFilesDir().getAbsolutePath();
    Mobile.setAndroidDataPath(filesDir);

    //Pasamos a golang el idioma del dispositivo.
    String langId = Locale.getDefault().getLanguage();
    Mobile.setAndroidLanguage(langId);

    new Thread(new Runnable() {
      @Override
      public void run() {
        Seq.setContext(getApplicationContext());
      }
    }).start();
  }

  @Override
  protected void onPause() {
    super.onPause();
    if(this.ebitenView != null) {
      this.ebitenView.suspendGame();
    }
  }

  @Override
  protected void onResume() {
    super.onResume();
    // Es buena práctica re-aplicar el modo inmersivo en onResume
    // ya que puede perderse si la app pierde el foco temporalmente.
    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
        hideSystemBarsApi30();
    } else {
        hideSystemBarsLegacy();
    }

    if(this.ebitenView != null) {
      this.ebitenView.resumeGame();
    }
  }

  private EbitenView getEbitenView() {
      return (EbitenView)this.findViewById(R.id.ebitenview);
  }

  private int hideSystemBars() {
    int uiOptions = 
      View.SYSTEM_UI_FLAG_LAYOUT_STABLE
    | View.SYSTEM_UI_FLAG_IMMERSIVE_STICKY        // Este es el flag clave para el modo "pegajoso"
    | View.SYSTEM_UI_FLAG_FULLSCREEN
    | View.SYSTEM_UI_FLAG_LAYOUT_FULLSCREEN
    | View.SYSTEM_UI_FLAG_HIDE_NAVIGATION         // Oculta la barra de navegación
    | View.SYSTEM_UI_FLAG_LAYOUT_HIDE_NAVIGATION; // Oculta la barra de estado
    return uiOptions;
  }

  //----------------------------------------------------------------------------------------------
  // Métodos para API 30+ (Android 11 y superior)
  //----------------------------------------------------------------------------------------------
  private void hideSystemBarsApi30() {
    // Get the WindowInsetsControllerCompat instance.
    // It handles the underlying platform API differences.
    WindowInsetsControllerCompat insetsController = WindowCompat.getInsetsController(getWindow(), getWindow().getDecorView());
    if (insetsController == null) {
      return;
    }

    // 1. Set the behavior for how bars reappear (sticky immersive)
    insetsController.setSystemBarsBehavior(WindowInsetsControllerCompat.BEHAVIOR_SHOW_TRANSIENT_BARS_BY_SWIPE);

    // 2. Hide the system bars using the COMPAT type
    insetsController.hide(WindowInsetsCompat.Type.systemBars());

    // 3. CONTROL THE ICON COLORS HERE!
    // To make status bar icons LIGHT (white) - false means light icons for a dark background
    insetsController.setAppearanceLightStatusBars(false);

    // To make navigation bar icons LIGHT (white) - false means light icons for a dark background
    insetsController.setAppearanceLightNavigationBars(false);

    // Optional: Listener for debugging/advanced control.
    insetsController.addOnControllableInsetsChangedListener(new WindowInsetsControllerCompat.OnControllableInsetsChangedListener() {
      @Override
      public void onControllableInsetsChanged(@NonNull WindowInsetsControllerCompat controller, int typeMask) {
        // Use WindowInsetsCompat.Type here as well
        if ((typeMask & WindowInsetsCompat.Type.systemBars()) != 0) {
          Log.d(TAG, "System bars became visible. Type mask: " + typeMask);
          // Re-apply appearance if needed (e.g., if system somehow reset them)
          controller.setAppearanceLightStatusBars(false);
          controller.setAppearanceLightNavigationBars(false);
        }
      }
    });
  }

  //----------------------------------------------------------------------------------------------
  // Métodos para APIs Legacy (menos de 30)
  //----------------------------------------------------------------------------------------------
  @SuppressWarnings("deprecation") // Suprimimos la advertencia de deprecación
  private void hideSystemBarsLegacy() {
    // Establece los flags de visibilidad del sistema
    View decorView = getWindow().getDecorView();
    decorView.setSystemUiVisibility(hideSystemBars());

    // Escuchar los cambios en la visibilidad de la UI para re-aplicar los flags
    // si se pierden (ej. al aparecer el teclado).
    decorView.setOnSystemUiVisibilityChangeListener(new View.OnSystemUiVisibilityChangeListener() {
      @Override
      public void onSystemUiVisibilityChange(int visibility) {
        if ((visibility & View.SYSTEM_UI_FLAG_FULLSCREEN) == 0) {
          // Si el flag FULLSCREEN se ha limpiado, significa que la barra
          // de estado es visible. Re-aplicamos los flags.
          decorView.setSystemUiVisibility(hideSystemBars());
        }
      }
    });
  }

  private EbitenView ebitenView;
  private View decorView;
}
