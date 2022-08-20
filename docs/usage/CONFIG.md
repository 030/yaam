# Config

```bash
mkdir ~/.yaam
chown 9999 -R ~/.yaam/
```

## Gradle

Adjust the `build.gradle` and/or `settings.gradle`:

```bash
repositories {
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/releases/'
    authentication {
      basic(BasicAuthentication)
    }
    credentials {
      username "hello"
      password "world"
    }
  }
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/3rdparty-maven/'
    authentication {
      basic(BasicAuthentication)
    }
    credentials {
      username "hello"
      password "world"
    }
  }
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/3rdparty-maven-gradle-plugins/'
    authentication {
      basic(BasicAuthentication)
    }
    credentials {
      username "hello"
      password "world"
    }
  }
}
```
