import {Client4, ClientError} from '@mattermost/client';
import {Options, ClientResponse} from '@mattermost/types/client4';

import {pluginId} from './manifest';
import {Customer, CustomerConfigValues, CustomerFilterOptions, CustomerPacketValues, CustomerPluginValues, FullCustomerInfo, GetCustomerResult} from './types/customers';

let siteURL = '';
let basePath = '';
let apiUrl = `${basePath}/plugins/${pluginId}/api/v0`;
const client = new Client4();
export const setSiteUrl = (url?: string): void => {
    if (url) {
        basePath = new URL(url).pathname.replace(/\/+$/, '');
        siteURL = url;
    } else {
        basePath = '';
        siteURL = '';
    }

    apiUrl = `${basePath}/plugins/${pluginId}/api/v0`;
};

export const getSiteUrl = (): string => {
    return siteURL;
};

export const getApiUrl = (): string => {
    return apiUrl;
};

export const doGet = async <TData = unknown>(url: string) => {
    const {data} = await doFetchWithResponse<TData>(url, {method: 'get'});

    return data;
};

export const doPost = async <TData = unknown>(url: string, body = {}) => {
    const {data} = await doFetchWithResponse<TData>(url, {
        method: 'POST',
        body,
    });

    return data;
};

export const doPut = async <TData = unknown>(url: string, body = {}) => {
    const {data} = await doFetchWithResponse<TData>(url, {
        method: 'PUT',
        body: JSON.stringify(body),
    });

    return data;
};

export const doFetchWithResponse = async <TData = unknown>(url: string, options: Options = {}): Promise<Omit<ClientResponse<TData | undefined>, 'headers'>> => {
    const response = await fetch(url, client.getOptions(options));
    let data;
    if (response.ok) {
        const contentType = response.headers.get('content-type');
        if (contentType === 'application/json') {
            data = await response.json() as TData;
        }

        return {
            response,
            data,
        };
    }

    const text = await response.text();
    throw new ClientError(client.url, {
        message: text || '',
        status_code: response.status,
        url,
    });
};

export function clientFetchCustomerByID(customerID: string) {
    return doGet<FullCustomerInfo>(`${apiUrl}/customers/${customerID}`);
}

export function clientFetchCustomers(opts?: CustomerFilterOptions) {
    const params = new URLSearchParams(opts).toString();
    return doGet<GetCustomerResult>(`${apiUrl}/customers${params ? `?${params}` : ''}`);
}

export function updateCustomer(customerID: string, customer: Partial<Customer>) {
    return doPut<FullCustomerInfo>(`${apiUrl}/customers/${customerID}`, customer);
}

export function updateCustomerConfig(customerID: string, config: Partial<CustomerConfigValues>) {
    return doPut<FullCustomerInfo>(`${apiUrl}/customers/${customerID}/config`, config);
}

export function updateCustomerPacket(customerID: string, packet: Partial<CustomerPacketValues>) {
    return doPut<FullCustomerInfo>(`${apiUrl}/customers/${customerID}/packet`, packet);
}

export function updateCustomerPlugins(customerID: string, plugins: Partial<CustomerPluginValues>) {
    return doPut<FullCustomerInfo>(`${apiUrl}/customers/${customerID}/plugins`, plugins);
}
