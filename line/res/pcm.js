//
// pcm.js pulse code modulated audio capture and replay code.
//

// Licensed under Creative Commons CC0.
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

export const name = 'pcm';

class PcmTransmitter extends AudioWorkletProcessor {
    constructor () {
        super();
    }

    process (inputs, outputs, parameters) {
        const input = inputs[0];
        const output = outputs[0];
        const channel = 0
        const arr16bit = new Int16Array(input[channel].length)
        const length = input[channel].length
        for (let i = 0; i < length; i++) {
            output[channel][i] = input[channel][i]
            arr16bit[i] = input[channel][i] * 32768
        }

        this.port.postMessage({
            eventType: 'data',
            audioBuffer: arr16bit
        });
        return true
    }
}

registerProcessor('pcm_send', PcmTransmitter);

class PcmReceiver extends AudioWorkletProcessor {
    constructor () {
        super();
        // It is acceptable to overflow.
        this.port.buf = new Float32Array(65536)
        this.port.begin = 0
        this.port.end = 0
        this.port.onmessage = function (ev) {
            const arr16bit = ev.data.audioBuffer;
            const length = arr16bit.length
            for (let i = 0; i < arr16bit.length; i++) {
                this.buf[(this.end + i) % 65536] = arr16bit[i] / 32768
            }
            this.end = (this.end + length) % 65536
        }
    }

    process (inputs, outputs, parameters) {
        const output = outputs[0]
        const channel = 0
        let buffer = this.port.buf
        let begin = this.port.begin
        let end = this.port.end

        for (let i = 0; i < output[channel].length; i++) {
            if (begin !== end) {
                output[channel][i] = buffer[begin]
                buffer[begin] = 0
                begin = (begin + 1) % 65536
            }
        }

        this.port.begin = begin
        return true
    }
}

registerProcessor('pcm_recv', PcmReceiver)
