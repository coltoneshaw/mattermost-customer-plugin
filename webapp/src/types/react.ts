import {Dispatch, SetStateAction} from 'react';

export type SetStateDispatch<T> = Dispatch<SetStateAction<T>>;
