import {Dispatcher} from "./dispatcher.js";
import {ApiClient} from "./api.js";
import {GetE, NewEC, NewECT} from "./util.js";

// Holds page-wide services; each one is created on first use.
export class App {
    /** @type {Dispatcher|null} */
    dispatcher = null;
    /** @type {ApiClient|null} */
    apiClient = null;

    /**
     * @returns {Dispatcher}
     */
    getDispatcher() {
        if (this.dispatcher === null) {
            this.dispatcher = new Dispatcher();
            this.dispatcher.addListener('error', (error) => this.showError(error));
        }

        return this.dispatcher;
    }

    /**
     * @returns {ApiClient}
     */
    getApiClient() {
        if (this.apiClient === null) {
            this.apiClient = new ApiClient(this.getDispatcher());
        }

        return this.apiClient;
    }

    /**
     * Appends a closable error message to the page's #errors div.
     * @param {{message: string}} error
     */
    showError(error) {
        console.error("error", error);

        const errors = GetE("errors");
        if (errors === null) {
            console.error("errors div not found");
            return;
        }

        const errorDiv = NewEC('div', 'error');
        errorDiv.appendChild(NewECT('span', 'error-message', error.message));

        const closeButton = NewECT('button', 'error-close', '✕');
        closeButton.addEventListener('click', () => errorDiv.remove());
        errorDiv.appendChild(closeButton);

        errors.appendChild(errorDiv);
    }
}
