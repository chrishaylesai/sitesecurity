import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Navigation } from "@/components/Navigation";
import { AuthContextProvider } from "@/lib/auth";

const inter = Inter({
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "SiteSecurity",
  description: "Security workforce management",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en-GB">
      <body className={`${inter.className} bg-gray-50`}>
        <AuthContextProvider>
          <Navigation />
          <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
            {children}
          </main>
        </AuthContextProvider>
      </body>
    </html>
  );
}
