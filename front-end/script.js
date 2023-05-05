$(document).ready(function() {
	// Carrega a lista de pessoas ao carregar a página
	loadPeopleList();

	// Manipula o envio do formulário de cadastro de pessoa
	$("#form-add").submit(function(event) {
		event.preventDefault();
		addPerson();
	});

	// Manipula a exclusão de uma pessoa da lista
	$("#table-body").on("click", ".btn-delete", function() {
		var id = $(this).data("id");
		deletePerson(id);
	});
});

function loadPeopleList() {
	$.ajax({
		url: "http://localhost:8080/pessoa",
		method: "GET",
		success: function(data) {
			$("#table-body").empty();
			data.forEach(function(person) {
				var row = $("<tr>");
				row.append($("<td>").text(person.id));
				row.append($("<td>").text(person.nome));
				row.append($("<td>").text(person.cpf));
				row.append($("<td>").text(person.endereco));
				row.append($("<td>").html('<button class="btn btn-danger btn-delete" data-id="' + person.id + '">Excluir</button>'));
				$("#table-body").append(row);
			});
		}
	});
}

function addPerson() {
    var form = document.getElementById("form-add");
    var formData = new FormData(form);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:8080/pessoa");
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onload = function() {
        if (xhr.status === 201) {
            alert("Pessoa adicionada com sucesso!");
            loadPeopleList()
            form.reset();
        } else {
            alert("Erro ao adicionar pessoa.");
        }
    };
    xhr.send(JSON.stringify(Object.fromEntries(formData.entries())));
}

/*document.getElementById("form-add").addEventListener("submit", function(event) {

});*/

function deletePerson(id) {
	$.ajax({
		url: "http://localhost:8080/pessoa/" + id,
		method: "DELETE",
		success: function() {
			loadPeopleList();
		}
	});
}
