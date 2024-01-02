import {Client4} from 'mattermost-redux/client';
import {Options, ClientResponse} from '@mattermost/types/client4';

import {pluginId} from './manifest';
import {CustomerFilterOptions, FullCustomerInfo, GetCustomerResult} from './types/customers';

let siteURL = '';
let basePath = '';
let apiUrl = `${basePath}/plugins/${pluginId}/api/v0`;

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

export const doFetchWithResponse = async <TData = unknown>(url: string, options: Options = {}): Promise<Omit<ClientResponse<TData | undefined>, 'headers'>> => {
    const response = await fetch(url, Client4.getOptions(options));
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

    data = await response.text();
    throw new Error(data, {
        cause: 'status code:' + response.status + '. url:' + url,
    });
};

export function clientFetchCustomerByID(customerID: string) {
    return doGet<FullCustomerInfo>(`${apiUrl}/customers/${customerID}`);
}

export function clientFetchCustomers(opts?: CustomerFilterOptions) {
    const params = new URLSearchParams(opts).toString();
    return doGet<GetCustomerResult>(`${apiUrl}/customers${params ? `?${params}` : ''}`);
}
