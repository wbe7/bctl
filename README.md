# bctl

## Использование

#### Инициализация проекта
```bash
bctl init --project-name webapp

# Помощь
bctl init --help
```

#### Добавление модуля
```bash
bctl add frontend --project-name webapp --module-image web/webapp --module-version 1.2.3 --module-port 8181 --ingress-class nginx

# Помощь
bctl add --help
```

#### Найти последнюю версию утилиты можно по [ссылке](https://git.tccenter.ru/tc-center/infra/App/base/bctl/-/releases).