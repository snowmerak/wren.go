import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
    // Get LSP server path from config or use default
    const config = workspace.getConfiguration('wren.lsp');
    let serverPath = config.get<string>('serverPath', '');
    
    if (!serverPath) {
        // Default to wren-lsp-std in PATH
        serverPath = 'wren-lsp-std';
    }

    const serverOptions: ServerOptions = {
        command: serverPath,
        args: []
    };

    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'wren' }],
        synchronize: {
            fileEvents: workspace.createFileSystemWatcher('**/*.wren')
        }
    };

    client = new LanguageClient(
        'wrenLanguageServer',
        'Wren Language Server',
        serverOptions,
        clientOptions
    );

    client.start();
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
