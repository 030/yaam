# Maven

~/.yaam/conf/caches.yaml

```bash
---
mavenReposAndUrls:
  3rdparty-maven: https://repo.maven.apache.org/maven2/
  3rdparty-maven-gradle-plugins: https://plugins.gradle.org/m2/
  3rdparty-maven-spring: https://repo.spring.io/release/
```

~/.yaam/conf/repositories/maven.yaml

```bash
---
allowedRepos:
  - releases
```

~/.yaam/conf/groups.yaml

```bash
---
groups:
  hello:
    - maven/releases
    - maven/3rdparty-maven
    - maven/3rdparty-maven-gradle-plugins
    - maven/3rdparty-maven-spring
```

## Gradle

### Preserve

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

### publish

```bash
publishing {
  publications {
    mavenJava(MavenPublication) {
      versionMapping {
        usage('java-api') {
          fromResolutionOf('runtimeClasspath')
        }
        usage('java-runtime') {
          fromResolutionResult()
        }
      }
    }
  }

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
  }
}
```
