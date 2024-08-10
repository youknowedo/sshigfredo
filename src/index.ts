import { timingSafeEqual } from "crypto";
import { readFileSync } from "fs";
import { inspect } from "util";

import ssh2 from "ssh2";
const { Server } = ssh2;

const server = new Server(
    {
        hostKeys: [readFileSync("host.key")],
    },
    (client) => {
        console.log("Client connected!");

        client
            .on("authentication", (ctx) => {
                ctx.accept();
            })
            .on("ready", () => {
                console.log("Client authenticated!");

                client.on("session", (accept, reject) => {
                    const session = accept();
                    session.once("exec", (accept, reject, info) => {
                        console.log(
                            "Client wants to execute: " + inspect(info.command)
                        );
                        const stream = accept();
                        stream.stderr.write("Oh no, the dreaded errors!\n");
                        stream.write("Just kidding about the errors!\n");
                        stream.exit(0);
                        stream.end();
                    });
                });
            })
            .on("close", () => {
                console.log("Client disconnected");
            });
    }
).listen(21, "127.0.0.1", () => {
    console.log("Listening on port " + (server.address() as any).port);
});
