import * as Sentry from "@sentry/browser";
import { BrowserTracing } from "@sentry/tracing";

const SENTRY_DSN = process.env.REACT_APP_SENTRY_DSN || "";

export const initSentry = () => {
  Sentry.init({
    dsn: SENTRY_DSN,
    integrations: [
      new BrowserTracing({
        tracePropagationTargets: ["localhost", "phenotype.dev", /^\/api\//],
      }),
    ],
    environment: process.env.NODE_ENV || "development",
    release: process.env.REACT_APP_VERSION || "development",
    tracesSampleRate: 0.1,
    replaysSessionSampleRate: 0.1,
    replaysOnErrorSampleRate: 1.0,
    ignoreErrors: [
      "ResizeObserver loop limit exceeded",
      "ResizeObserver loop completed with undelivered notifications",
      "Network Error",
    ],
    beforeSend(event, hint) {
      // Filter out non-critical errors in development
      if (process.env.NODE_ENV === "development") {
        const error = hint.originalException;
        if (error && error.message && error.message.includes("Warning:")) {
          return null;
        }
      }
      return event;
    },
  });
};

export const setUserContext = (user: { id: string; email?: string; username?: string }) => {
  Sentry.setUser({
    id: user.id,
    email: user.email,
    username: user.username,
  });
};

export const clearUserContext = () => {
  Sentry.setUser(null);
};

export const captureMessage = (message: string, level: Sentry.SeverityLevel = "info") => {
  Sentry.captureMessage(message, level);
};

export const captureException = (error: Error, context?: Record<string, unknown>) => {
  if (context) {
    Sentry.captureException(error, { extra: context });
  } else {
    Sentry.captureException(error);
  }
};

export const addBreadcrumb = (
  category: string,
  message: string,
  level: Sentry.SeverityLevel = "info"
) => {
  Sentry.addBreadcrumb({
    category,
    message,
    level,
    timestamp: Date.now() / 1000,
  });
};

export default {
  initSentry,
  setUserContext,
  clearUserContext,
  captureMessage,
  captureException,
  addBreadcrumb,
};
