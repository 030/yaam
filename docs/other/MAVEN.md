# Maven

## Gradle

Adjust the repositories sections in the build.gradle and settings.gradle:

```bash
repositories {
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/maven/groups/hello/'
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
