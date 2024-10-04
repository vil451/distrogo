import http from 'node:http';

// Проверка версии Node.js
const version = process.versions.node;

const requestHandler = (req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });

    // Используем переменную для хранения вывода
    let responseText = `Hello World! Running on Node.js version: ${version}\n`;

    // Проверяем, доступен ли import.meta.resolve (для Node.js 20+)
    if (typeof import.meta.resolve === 'function') {
        try {
            const resolvedPath = import.meta.resolve('fs');
            responseText += `Resolved path for 'fs' module: ${resolvedPath}\n`;
        } catch (err) {
            responseText += 'import.meta.resolve is not supported in this version of Node.js\n';
        }
    } else {
        // Если функция недоступна (Node.js 18 или более ранняя версия)
        responseText += 'import.meta.resolve is not available in this version of Node.js\n';
    }

    // Завершаем запись в ответе и отправляем его
    res.end(responseText);
};

// Создаем HTTP сервер
const server = http.createServer(requestHandler);

// Сервер слушает на порту 3000
server.listen(3000, () => {
    console.log('Server is running on http://localhost:3000');
});
