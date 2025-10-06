"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode_1 = require("vscode");
const node_1 = require("vscode-languageclient/node");
let client;
let wrenTerminal;
function activate(context) {
    // Get LSP server path from config or use default
    const config = vscode_1.workspace.getConfiguration('wren.lsp');
    let serverPath = config.get('serverPath', '');
    if (!serverPath) {
        // Default to wren-lsp-std in PATH
        serverPath = 'wren-lsp-std';
    }
    const serverOptions = {
        command: serverPath,
        args: []
    };
    const clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'wren' }],
        synchronize: {
            fileEvents: vscode_1.workspace.createFileSystemWatcher('**/*.wren')
        }
    };
    client = new node_1.LanguageClient('wrenLanguageServer', 'Wren Language Server', serverOptions, clientOptions);
    client.start();
    // Register CLI commands
    context.subscriptions.push(vscode_1.commands.registerCommand('wren.runCurrentFile', async () => {
        const editor = vscode_1.window.activeTextEditor;
        if (!editor) {
            vscode_1.window.showErrorMessage('No active editor');
            return;
        }
        const document = editor.document;
        if (document.languageId !== 'wren') {
            vscode_1.window.showErrorMessage('Current file is not a Wren file');
            return;
        }
        await runWrenFile(document.uri);
    }));
    context.subscriptions.push(vscode_1.commands.registerCommand('wren.runFile', async (uri) => {
        if (!uri) {
            vscode_1.window.showErrorMessage('No file selected');
            return;
        }
        await runWrenFile(uri);
    }));
}
exports.activate = activate;
async function runWrenFile(fileUri) {
    const config = vscode_1.workspace.getConfiguration('wren.cli');
    let cliPath = config.get('path', '');
    if (!cliPath) {
        // Try common CLI names
        cliPath = 'wren-cli-std';
    }
    const filePath = fileUri.fsPath;
    // Get or create terminal
    if (!wrenTerminal || wrenTerminal.exitStatus !== undefined) {
        wrenTerminal = vscode_1.window.createTerminal('Wren');
    }
    wrenTerminal.show();
    wrenTerminal.sendText(`${cliPath} "${filePath}"`);
}
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map