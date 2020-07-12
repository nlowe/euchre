function roomCode() {
    const codeInput = document.getElementById('join-code');
    return codeInput.value !== '' ? codeInput.value : codeInput.placeholder;
}

function joinRoom() {
    window.location.href += (window.location.href.endsWith('/') ? 'table/' : '/table/') + roomCode() + '/';
}
