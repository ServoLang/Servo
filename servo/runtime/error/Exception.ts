/*
 * Copyright (c) 2024. Servo Contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

export enum ErrorType {
    SyntaxError = 1,
    TypeError = 1,
    RuntimeError = 1,
    ReferenceError = 1,
    InternalError = 1,
    DivideByZero = 1,
    UnknownError = 1,

    InvalidAssignment = 1,
    Finished = 0,
}

export class Exception extends Error {
    private readonly type: ErrorType;
    private readonly exitCode: number = 1;

    constructor (type: ErrorType, message?: string) {
        super(message);
        this.type = type;
        this.exitCode = type.valueOf();
        this.exit();
    }

    public getType (): ErrorType {
        return this.type;
    }

    public getExitCode (): number {
        return this.exitCode;
    }

    public exit (): Exception {
        if (this.getType().valueOf() == 0) console.error(`Error: ${this.message}`);
        else console.log(`Finished with exit code: (${this.getExitCode()})`);
        Deno.exit(this.getExitCode());
    }

    public override toString (): string {
        return `${this.message}`;
    }

}
