#!/usr/bin/env python

import socket
import time
import audioop
import math
import pyaudio
import typer
import json

ac = 0
samples = []
avg = 0
max_db = 0


def ag_samples(sample):
    global ac
    global samples
    global avg
    global max_db

    if ac < 1:
        if sample > 1:
            samples.append(sample)
            ac = ac+1
    else:
        avg = sum(samples) / len(samples)
        if float(avg) > float(max_db):
            max_db = "%.2f" % avg
        ac = 0
        samples = []

    return "%.2f" % avg


def clear_all():
    ''' Clear all the variables '''
    print("Inside Clear Function")
    global sound_tracks
    global max_value
    global max_db
    global current_selection
    max_value = [0] * 3
    max_db = 0
    current_selection = 0


def main(id: str = "abc"):
    ''' Main function '''
    CHUNK = 1024 * 4
    FORMAT = pyaudio.paInt16
    CHANNELS = 1
    RATE = 44100
    p = pyaudio.PyAudio()
    stream = p.open(format=FORMAT,
                    channels=CHANNELS,
                    rate=RATE,
                    input=True,
                    output=False,
                    frames_per_buffer=CHUNK)

    SERVER_IP = '127.0.0.1'
    SERVER_PORT = 8080

    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect((SERVER_IP, SERVER_PORT))
    while True:
        total = 0
        data = stream.read(CHUNK,
                           exception_on_overflow=False)
        reading = audioop.max(data, 2)

        total = 3*(math.log10(abs(reading)))

        db = ag_samples(total)
        client_socket.send(json.dumps({
            "id": id,
            "dB": db
        }).encode())
        time.sleep(.005)


if __name__ == '__main__':
    typer.run(main)
