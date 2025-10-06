"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode_1 = require("vscode");
const node_1 = require("vscode-languageclient/node");
let client;
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
}
exports.activate = activate;
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map