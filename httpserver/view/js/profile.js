function readURL(input) {
    if (input.files && input.files[0]) {
        const reader = new FileReader();
        reader.onload = function (e) {
            $('#imageResult').attr('src', e.target.result);
        };
        reader.readAsDataURL(input.files[0]);
        enableSubmitButton()
    }
}

function onNicknameChange() {
    const nickname = document.getElementById("nickname")
    const originalNickname = nickname.defaultValue
    console.log(originalNickname)
    console.log(nickname.value)

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
