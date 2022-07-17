# YAAM

Yet Another Artifact Manager

## Start

```bash
./yaam
```

## Configure

Adjust the `build.gradle`:

```bash
repositories {
    maven {
        allowInsecureProtocol true
        url 'http://localhost:25113/maven2central/'
    }
}
```

## Cache artifacts

```bash
./gradlew b
```
