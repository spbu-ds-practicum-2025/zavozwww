document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem("token");
    
    if (token) {
        document.getElementById("navbar").classList.remove("hidden");
        const app = new App();
    } else {
        const authManager = new AuthManager();
    }
    console.log('App start');
});