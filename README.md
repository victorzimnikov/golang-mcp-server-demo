# Демо MCP-сервер на Golang

Демо проект MCP-сервера на Golang. Предоставляет инструменты для Claude и Codex:

- `create_task`
- `get_project_context`
- `list_projects`
- `server_info`
- `update_task_status`

## Скрипты

- `make run_inspector` - запуск mcp-инспектора
- `make build_server` - сборка бинарника
- `make dev_server` - запуск локального dev-сервера

## Подключение

### Claude

```
claude mcp add --transport http mcp-test-hub --scope local http://127.0.0.1:8000/mcp
```

### Codex

```
codex mcp add mcp-test-hub --url http://127.0.0.1:8000/mcp
```

## Инструменты

### `create_task` - Создание задачи

Промпт:

```
Используя mcp-test-hub, создай для проекта #1 задачу "Test task" с описанием "Test task description" и низким приоритетом
```

Входные данные:

```
{
  "project_id": 1,
  "title": "Test task",
  "description": "Test task description",
  "priority": "low"
}
```

Выходные данные:

```
{
  "project_id": 1,
  "title": "Test task",
  "description": "Test task description",
  "priority": "low",
  "status": "todo",
  "source": "claude",
  "version": 1
}
```

### `list_projects` - Получение списка проектов

Промпт:

```
Используя mcp-test-hub, покажи список проектов
```

Входные данные:

```
{}
```

Выходные данные:

```
{
  "projects": [
    {
      "created_at": "2026-09-28T13:29:55Z",
      "description": "Shared context demo",
      "id": 1,
      "name": "Project Context Hub",
      "updated_at": "2026-09-28T13:29:55Z"
    }
  ]
}
```

### `server_info` - Получение информации о сервере

Промпт:

```
Используя mcp-test-hub, покажи информацию о сервере
```

Входные данные:

```
{}
```

Выходные данные:

```
{
  "name": "golang-mcp-server-demo",
  "version": "v0.1.0"
}
```

### `update_task_status` - Изменение статуса задачи

Промпт:

```
Используя mcp-test-hub, переведи задачу #3 в статус в работе
```

Входные данные:

```
{
  "task_id": 3,
  "status": "in_progress",
  "expected_version": 1
}
```

Выходные данные:

```
{
  "created_at": "2026-09-28T21:43:37Z",
  "description": "Реализовать фичу....",
  "id": 3,
  "priority": "medium",
  "project_id": 1,
  "source": "claude",
  "status": "in_progress",
  "title": "Начать реализацию",
  "updated_at": "2026-09-29T14:01:17Z",
  "version": 2
}
```

### `get_project_context` - Получение контекста проекта

Промпт:

```
Используя mcp-test-hub, покажи контекст проекта #1
```

Входные данные:

```
{
  "project_id": 3
}
```

Выходные данные:

```
{
  "open_tasks": [
    {
      "created_at": "2026-09-28T21:20:06Z",
      "description": "",
      "id": 2,
      "priority": "low",
      "project_id": 1,
      "source": "claude",
      "status": "in_progress",
      "title": "111",
      "updated_at": "2026-09-29T11:40:57Z",
      "version": 2
    },
    {
      "created_at": "2026-09-28T21:43:37Z",
      "description": "Реализовать фичу....",
      "id": 3,
      "priority": "medium",
      "project_id": 1,
      "source": "claude",
      "status": "in_progress",
      "title": "Начать реализацию",
      "updated_at": "2026-09-29T14:01:17Z",
      "version": 2
    },
    {
      "created_at": "2026-09-29T05:47:25Z",
      "description": "",
      "id": 4,
      "priority": "medium",
      "project_id": 1,
      "source": "codex",
      "status": "todo",
      "title": "Проверить подключение Codex",
      "updated_at": "2026-09-29T05:47:25Z",
      "version": 1
    },
    {
      "created_at": "2026-09-29T13:49:39Z",
      "description": "Test task description",
      "id": 5,
      "priority": "low",
      "project_id": 1,
      "source": "claude",
      "status": "todo",
      "title": "Test task",
      "updated_at": "2026-09-29T13:49:39Z",
      "version": 1
    }
  ],
  "project": {
    "created_at": "2026-09-28T13:29:55Z",
    "description": "Shared context demo",
    "id": 1,
    "name": "Project Context Hub",
    "updated_at": "2026-09-28T13:29:55Z"
  }
}
```
