//
// AES symmetric key encryption.
//

// Licensed under Creative Commons CC0.
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// The counter block value must never be reused with a given key.
// counter = window.crypto.getRandomValues(new Uint8Array(16))
//let counter = new Uint8Array(16) // TODO window.crypto.getRandomValues(new Uint8Array(16))

export { generateKey, importKey, decrypt, encrypt, redirectWithSecrets };

async function generateKey(url) {
    return await window.crypto.subtle.generateKey(
        {
            name: "AES-CTR",
            length: 256
        },
        true,
        ["encrypt", "decrypt"]
    ).then(async (newkey) => {
        return await window.crypto.subtle.exportKey("jwk", newkey).then((exp) => {
            let ret = 'You can join the call with the following url ' + url + '#leaf_' + exp.k + ' .'
            console.log(ret)
            return ret
        })
    })
}

async function redirectWithSecrets() {
    window.crypto.subtle.generateKey(
        {
            name: "AES-CTR",
            length: 64
        },
        true,
        ["encrypt", "decrypt"]
    ).then(async (key) => {
        return await window.crypto.subtle.exportKey("jwk", key)
    }).then((exp) => '#leaf_' + exp.k
    ).then(leaf => {
        let apikey = btoa(window.crypto.getRandomValues(new Uint8Array(16)).toString())
        apikey = apikey.replaceAll('=', '0')
        window.location.href = document.location.origin + document.location.pathname + "?apikey=" + apikey + leaf
    })
}

function getKeyFromMessage(key1) {
    let i = key1.indexOf('#')
    if (i !== -1) {
        key1 = key1.substring(i + 1, key1.length)
    }
    if (key1.startsWith("leaf_")) {
        key1 = key1.substring(5, key1.length)
    }
    i = key1.indexOf(' ')
    if (i !== -1) {
        key1 = key1.substring(0, i)
    }
    if (key1.length > 0) {
        return new Promise((res, _) => {
            res(key1)
        })
    } else {
        return new Promise((_, rej) => {
            rej(key1)
        })
    }
}

async function importKey(key1) {
    return await getKeyFromMessage(key1
    ).then(async key2 => {
            let aes = JSON.stringify({"alg":"A256CTR","ext":true,"k":key2,"key_ops":["encrypt","decrypt"],"kty":"oct","ws":"wss://line.eper.io"})
            return JSON.parse(aes)
    }).then(async (key3) => {
        return await window.crypto.subtle.importKey(
            "jwk",
            key3,
            "AES-CTR",
            true,
            ["encrypt", "decrypt"])
    })
}

async function encrypt(cypher, counter, arrayOfWords) {
    return await window.crypto.subtle.encrypt(
        {
            name: "AES-CTR",
            counter,
            length: 64
        },
        cypher,
        arrayOfWords)
        .then((ciphertext) => {
            return new Int16Array(ciphertext, 0, arrayOfWords.length);
        });
}

async function decrypt(cypher, counter, arrayOfWords) {
    return await window.crypto.subtle.decrypt(
        {
            name: "AES-CTR",
            counter,
            length: 64
        },
        cypher,
        arrayOfWords)
        .then((decrypted) => {
            return new Int16Array(decrypted);
        });
}

// Reference: https://github.com/mdn/dom-examples/blob/master/web-crypto/encrypt-decrypt/aes-ctr.js
