class TemporaryNotice {
    show(message, type){
        const noticeContainer = document.createElement("div");
        noticeContainer.classList.add("temp-notice");
        noticeContainer.classList.add(`${type}`);
        noticeContainer.innerHTML = `
            <p class=temp-notice__message>${message}</p>
        `

        document.body.appendChild(noticeContainer);

        setTimeout(() => {
            noticeContainer.classList.add("show");
        }, 10);

        setTimeout(() => {
            noticeContainer.classList.remove("show");
            noticeContainer.classList.add("hide");

            setTimeout(() => {
                document.body.removeChild(noticeContainer);
            }, 300);
        }, 3000);
    }

    success(message){
        return this.show(message, "success");
    }

    error(message){
        return this.show(message, "error");
    }
}

const tempNotice = new TemporaryNotice();