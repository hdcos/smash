package io.smash.interpreter;

import java.io.IOException;

import io.smash.parser.AST;

public class Interpreter {

    public static void exec(AST root) throws IOException {

        if (root == null) {
            return;
        }

        switch (root.which()) {

            case CMD: {
                var cmd = root.command();
                var process = Runtime.getRuntime().exec(cmd.toArray(new String[cmd.size()]));

                try {
                    var res = process.waitFor();
                    if (res != 0) {
                        System.err.println("exec failed");
                    }
                } catch (InterruptedException e) {
                    process.destroy();
                }
                var reader = process.inputReader();
                String read = reader.readLine();
                while (read != null) {
                    System.out.println(read);
                    read = reader.readLine();
                }
                break;
            }

        }

    }
}
