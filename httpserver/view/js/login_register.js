function revealPassword(e) {
    e.stopPropagation()
    const passwordField = document.getElementById("password");
    const toggle = document.getElementById("togglePassword")
    if (passwordField.type === "password") {
        passwordField.type = "text";
        toggle.classList.remove("fa-eye-slash")
        toggle.classList.add("fa-eye")
    } else {
        passwordField.type = "password";
        toggle.classList.add("fa-eye-slash")
        toggle.classList.remove("fa-eye")
    }
}