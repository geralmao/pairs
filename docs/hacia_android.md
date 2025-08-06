# Hacia dispositivos móviles
Para generar un producto hacia un dispositivo móvil, **Ebiten** nos proporciona el paquete **ebitengine/gomobile** que trabaja junto con herramientas como **ebitenmobile**.

Nos va a permitir generar librerías nativas tanto para **Android** como para **iOS**:
* .aar (Android Archive) en Android
* .framework en iOS.

Nos permite usar **Ebiten** a través de aplicaciones nativas **Java/Kotlin** y **Swift** respectivamente.

## Instalación de herrameintas.
Debemos instalar **ebitenmobile**:

~~~shell
go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest
~~~

Una vez instaladas, se podrá generar los **bindings** nativos para Android e iOS:

### Android.
Para generar la librería compartida:

~~~shell
ebitenmobile bind -target android -javapkg <com.yourcompany.paquete> -o <nombrelibreria.aar>
~~~

Si no encuentra **ebitenmobile** hay que añadir la ruta de busqueda de la siguiente forma:
~~~shell
export PATH=$PATH:$(go env GOPATH)/bin
~~~

* Con esto se complia el código **Go** en una librería compartida **.so** para Android.
* Empaqueta el código nativo y los bindings Java en un archivo **.aar**

#### Las clases bindings.
Son las clases puente que se generan a partir del código **Go** con funciones exportadas con el comentario especial o directiva `//export NombreFunción`. 
Esta directiva sólo se puede usar en:
* funciones globales (no puede ser aplicado en receptores (receivers))
* en funciones globales del paquete principal o desde otros paquetes siempre que sea visible.

Las clases puente se guardan en un fichero llamado **clasess.jar** y bajo **jni** se guardarán las librerías compartidas por plataforma y como hemos indicado todo esto dentro del fichero **.aar**.

#### Estructura para Android.
Una vez que tengamos el fichero **.aar**, ese lo ubicaremos en un proyecto Android básico:
~~~bash
android/
├── app/
│   ├── build.gradle
│   ├── libs/
│   │   └── nombrelibreria.aar
│   └── src/
│       └── main/
│           ├── AndroidManifest.xml
│           └── java/
│               └── com/
│                   └── example/
│                       └── nombreproyecto/
│                           └── MainActivity.java
├── build.gradle
└── settings.gradle
~~~

para ir creando la estructura del proyecto usamos **gradle**:
~~~shell
gradle init --type basic --dsl groovy --project-name "nombre proyecto"  --use-defaults
~~~

Donde el contenido de cada fichero puede ser:

##### settings.gradle (raiz)
~~~gradle
rootProject.name = 'nombreproyecto'
include ':app'
~~~

##### build.gradle (raiz)
~~~gradle
plugins {
  id 'com.android.application' version '8.2.0' apply false 
}

buildscript {
  ext {
    agp_version = '8.2.0' // Versión del Android Gradle Plugin
  }
  repositories {
    google()
    mavenCentral()
  }
  dependencies {
    classpath "com.android.tools.build:gradle:$agp_version"
  }
}

allprojects {
  repositories {
    google()
    mavenCentral()
  }
}
~~~

creamos el directorio **app** y sus subdirectorios:
~~~shell
mkdir -p app/libs && mkdir -p app/src/main/java
~~~


##### app/build.gradle
~~~gradle
plugins {
  id 'com.android.application'
}

android {
  namespace 'com.example.nombreproyecto' // Tu paquete base
  compileSdk 34 // La versión del SDK con la que compilas

  defaultConfig {
    applicationId "com.example.nombreproyecto" // ID de la aplicación
    minSdk 21 // Versión mínima de Android
    targetSdk 34 // Versión de destino
    versionCode 1
    versionName "1.0"

    testInstrumentationRunner "androidx.test.runner.AndroidJUnitRunner"
  }

  buildTypes {
    release {
      minifyEnabled false
      proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
    }
  }
  compileOptions {
    sourceCompatibility JavaVersion.VERSION_17 // O la versión de Java que estés usando
    targetCompatibility JavaVersion.VERSION_17
  }

  // Configuración para el AAR de Ebitenmobile
  packagingOptions {
    jniLibs {
      useLegacyPackaging = true // Esto es importante para el .so de Go
    }
  }
}

dependencies {
  // Dependencias estándar de Android
  implementation 'androidx.core:core-ktx:1.10.1'
  implementation 'androidx.appcompat:appcompat:1.6.1'
  implementation 'com.google.android.material:material:1.9.0'
  implementation 'androidx.constraintlayout:constraintlayout:2.1.4'
  testImplementation 'junit:junit:4.13.2'
  androidTestImplementation 'androidx.test.ext:junit:1.1.5'
  androidTestImplementation 'androidx.test.espresso:espresso-core:3.5.1'

  // ¡Aquí es donde añades tu AAR!
  implementation fileTree(dir: 'libs', include: ['*.aar'])
}
~~~

##### app/src/main/AndroidManifest.xml
~~~xml
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
  package="com.example.nombreproyecto">

  <application
    android:label="Nombre Proyecto"
    android:theme="@android:style/Theme.NoTitleBar.Fullscreen">
    <activity android:name=".MainActivity"
      android:configChanges="orientation|keyboardHidden|screenSize"
      android:label="Nombre Proyecto"
      android:screenOrientation="portrait"
      android:exported="true">
      <intent-filter>
        <action android:name="android.intent.action.MAIN" />
        <category android:name="android.intent.category.LAUNCHER" />
      </intent-filter>
    </activity>
  </application>
</manifest>
~~~

creamos el recurso de la vista que va a usar **Ebiten** bajo **res/layout** y la llamamos **activity_main.xml**:

~~~shell
mkdir -p app/src/main/res/layout
~~~

##### app/src/main/res/layout/activity_main.xml
~~~xml
<?xml version="1.0" encoding="utf-8"?>
<FrameLayout xmlns:android="http://schemas.android.com/apk/res/android"
  xmlns:tools="http://schemas.android.com/tools"
  android:id="@+id/ebitenview"
  android:layout_width="match_parent"
  android:layout_height="match_parent"
  tools:context=".MainActivity">
</FrameLayout>
~~~


creamos el paquete de la aplicación bajo **java** y añadimos la actividad de entrada:

~~~shell
mkdir -p app/src/main/java/com/example/nombreproyecto
~~~

##### app/src/main/java/com/example/nombreproyecto/MainActivity.java
~~~java
package com.example.nombreproyecto;

import android.app.Activity;
import android.os.Bundle;

public class MainActivity extends Activity {
  static {
    System.loadLibrary("gojni"); // nombre del .so generado por ebitenmobile
  }

  @Override
  protected void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    StartGame(); // Función exportada desde Go
  }

  public native void StartGame();
}
~~~

##### app/src/main/java/com/example/nombreproyecto/EbitenViewWithErrorHandling.java
~~~java
package com.example.nombreproyecto;

import android.content.Context;
import android.util.AttributeSet;
import com.programatta.go.pairs.ebitenmobileview.EbitenView;


class EbitenViewWithErrorHandling extends EbitenView {
  public EbitenViewWithErrorHandling(Context context) {
    super(context);
  }

  public EbitenViewWithErrorHandling(Context context, AttributeSet attributeSet) {
    super(context, attributeSet);
  }

  @Override
  protected void onErrorOnGameUpdate(Exception e) {
    // You can define your own error handling e.g., using Crashlytics.
    // e.g., Crashlytics.logException(e);
    super.onErrorOnGameUpdate(e);
  }
}
~~~

##### app/libs/nombrelibreria.aar
Copiamos el fichero generaro anteriormente a **app/libs**:

~~~shell
cp ../android-lib/nombrelibreria.aar app/libs
~~~


#### Generar el apk.
Para generar el **apk**, ya usaremos otro contenedor que contenga el **SDK de Android** con las herramientas necesarias.

En modo depuración:
~~~shell
./gradlew assembleDebug
~~~

En modo release para subir en tiendas:
~~~shell
./gradlew assembleRelease
~~~


### iOS.
Al igual que para generar el binario para **MacOS** se requerirá de una máquina con dicho sistema operativo y las herramientas asociadas.

