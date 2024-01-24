import {UseFormReturnType} from '@mantine/form';

export type FormTextInputParams<T> = {
    getInputProps: UseFormReturnType<T, (values: T) => T>['getInputProps'];
    label?: string;
    placeholder?: string;
    formKey: keyof T;
}
