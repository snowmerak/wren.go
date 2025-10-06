import * as path from 'path';
import { workspace, ExtensionContext, commands, window, Uri, Terminal } from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
} from 'vscode-languageclient/node';

let client: LanguageClient;
let wrenTerminal: Terminal | undefined;

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

    // Register CLI commands
    context.subscriptions.push(
        commands.registerCommand('wren.runCurrentFile', async () => {
            const editor = window.activeTextEditor;
            if (!editor) {
                window.showErrorMessage('No active editor');
                return;
            }
            
            const document = editor.document;
            if (document.languageId !== 'wren') {
                window.showErrorMessage('Current file is not a Wren file');
                return;
            }
            
            await runWrenFile(document.uri);
        })
    );

    context.subscriptions.push(
        commands.registerCommand('wren.runFile', async (uri: Uri) => {
            if (!uri) {
                window.showErrorMessage('No file selected');
                return;
            }
            
            await runWrenFile(uri);
        })
    );
}

async function runWrenFile(fileUri: Uri) {
    const config = workspace.getConfiguration('wren.cli');
    let cliPath = config.get<string>('path', '');
    
    if (!cliPath) {
        // Try common CLI names
        cliPath = 'wren-cli-std';
    }
    
    const filePath = fileUri.fsPath;
    
    // Get or create terminal
    if (!wrenTerminal || wrenTerminal.exitStatus !== undefined) {
        wrenTerminal = window.createTerminal('Wren');
    }
    
    wrenTerminal.show();
    wrenTerminal.sendText(`${cliPath} "${filePath}"`);
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
