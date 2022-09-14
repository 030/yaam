# Config

```bash
mkdir -p ~/.yaam/conf
chown 9999 -R ~/.yaam/
```

vim ~/.yaam/conf/caches.yaml

```bash
mavenReposAndUrls:
  3rdparty-maven: https://repo.maven.apache.org/maven2/
  3rdparty-maven-gradle-plugins: https://plugins.gradle.org/m2/
  3rdparty-maven-spring: https://repo.spring.io/release/
```

vim ~/.yaam/conf/groups.yaml

```bash
groups:
  hello:
    - releases
    - 3rdparty-maven
    - 3rdparty-maven-gradle-plugins
    - 3rdparty-maven-spring
```

vim ~/.yaam/conf/repositories.yaml

```bash
maven:
  - releases
```

## Gradle

Adjust the `build.gradle` and/or `settings.gradle`:

```bash
repositories {
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/maven/releases/'
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
    url 'http://localhost:25213/maven/3rdparty-maven/'
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
    url 'http://localhost:25213/maven/3rdparty-maven-gradle-plugins/'
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
