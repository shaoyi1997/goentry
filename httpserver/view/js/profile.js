function readURL(input) {
    const maxFileSize = 50000000

    if (input.files && input.files[0]) {
        if (input.files[0].size > maxFileSize) {
            input.value = "";
            document.getElementById("image-feedback").style.display = "inherit";
            return
        }

        const reader = new FileReader();
        reader.onload = function (e) {
            document.getElementById("image-feedback").src = e.target.result;
        };
        reader.readAsDataURL(input.files[0]);
        enableSubmitButton()
    }
}

function onClickUpload() {
    document.getElementById("image-feedback").style.display = "none";
}

function onNicknameChange() {
    const nickname = document.getElementById("nickname")
    const originalNickname = nickname.defaultValue

    if (nickname.value !== originalNickname) {
        enableSubmitButton()
    } else {
        disableSubmitButton()
    }
}

function disableSubmitButton() {
    document.getElementById("submit-button").classList.add("disabled");
}

function enableSubmitButton() {
    document.getElementById("submit-button").classList.remove("disabled");
}