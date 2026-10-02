const searchInput = document.getElementById("searchInput");
const searchButton = document.getElementById("searchButton");
const vacanciesContainer = document.getElementById("vacancies");
const status = document.getElementById("status");

async function searchVacancies() {
    const query = searchInput.value.trim();

    if (!query) {
        status.textContent = "Введите запрос";
        return;
    }

    status.textContent = "Ищем вакансии...";
    vacanciesContainer.innerHTML = "";

    try {
        const response = await fetch(
            `/api/vacancies?query=${encodeURIComponent(query)}`
        );

        if (!response.ok) {
            throw new Error("Ошибка сервера");
        }

        const vacancies = await response.json();

        status.textContent = `Найдено вакансий: ${vacancies.length}`;

        vacancies.forEach(vacancy => {
            const card = document.createElement("div");
            card.className = "vacancy";

            card.innerHTML = `
                <h2>${escapeHTML(vacancy.title)}</h2>

                <div class="company">
                    ${escapeHTML(vacancy.company)}
                </div>

                <div class="info">
                    <span class="salary">
                        ${escapeHTML(vacancy.salary)}
                    </span>

                    <span>
                        ${escapeHTML(vacancy.location)}
                    </span>
                </div>

                <a
                    href="${escapeAttribute(vacancy.url)}"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Открыть вакансию
                </a>
            `;

            vacanciesContainer.appendChild(card);
        });

    } catch (error) {
        status.textContent = "Не удалось получить вакансии";
        console.error(error);
    }
}

function escapeHTML(value) {
    const div = document.createElement("div");
    div.textContent = value ?? "";
    return div.innerHTML;
}

function escapeAttribute(value) {
    return String(value ?? "").replace(/"/g, "&quot;");
}

searchButton.addEventListener("click", searchVacancies);

searchInput.addEventListener("keydown", event => {
    if (event.key === "Enter") {
        searchVacancies();
    }
});