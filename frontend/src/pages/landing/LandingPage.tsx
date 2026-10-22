import { Link } from "react-router-dom"
import type { LucideIcon } from "lucide-react"
import {
  ArrowRight,
  Boxes,
  Eye,
  FileText,
  FolderTree,
  Image as ImageIcon,
  Link2,
  ShieldCheck,
  UploadCloud,
} from "lucide-react"

import { useAuth } from "../auth/application/useAuth"

interface Feature {
  tag: string
  icon: LucideIcon
  title: string
  description: string
}

const features: Feature[] = [
  {
    tag: "BUCKET",
    icon: FolderTree,
    title: "Buckets",
    description: "Group objects into buckets you can spin up or drop at will.",
  },
  {
    tag: "UPLOAD",
    icon: UploadCloud,
    title: "Upload anything",
    description: "Drop in images and documents with instant inline previews.",
  },
  {
    tag: "SHARE",
    icon: Link2,
    title: "Share links",
    description: "Mint a public link for any object in a single click.",
  },
  {
    tag: "ACCESS",
    icon: Eye,
    title: "Public or private",
    description: "Flip an object's visibility and expose a read URL on demand.",
  },
  {
    tag: "AUTH",
    icon: ShieldCheck,
    title: "Secure sessions",
    description: "JWT access tokens with rotating refresh tokens on every use.",
  },
]

const previewRows = [
  { name: "launch-hero.png", size: "412 KB", kind: "png", public: true },
  { name: "pricing-v3.pdf", size: "88 KB", kind: "pdf", public: false },
  { name: "notes.md", size: "3 KB", kind: "md", public: false },
]

export function LandingPage() {
  const { status } = useAuth()
  const authenticated = status === "authenticated"
  const cta = authenticated ? "/drive" : "/login"

  return (
    <div className="flex min-h-svh flex-col bg-paper">
      <header className="flex items-center justify-between border-b border-line px-6 py-4">
        <span className="flex items-center gap-2">
          <span className="flex h-7 w-7 items-center justify-center rounded-md bg-verdigris text-white">
            <Boxes className="h-4 w-4" />
          </span>
          <span className="font-display text-lg font-semibold tracking-tight text-ink">
            ObjexRG
          </span>
        </span>
        <nav className="flex items-center gap-2">
          {authenticated ? (
            <Link
              to="/drive"
              className="inline-flex items-center gap-2 rounded-lg bg-verdigris px-4 py-2 text-sm font-medium text-white hover:bg-verdigris-dark"
            >
              Open Drive
              <ArrowRight className="h-4 w-4" />
            </Link>
          ) : (
            <>
              <Link
                to="/login"
                className="rounded-lg px-4 py-2 text-sm font-medium text-graphite hover:text-ink"
              >
                Sign in
              </Link>
              <Link
                to="/login"
                className="inline-flex items-center gap-2 rounded-lg bg-verdigris px-4 py-2 text-sm font-medium text-white hover:bg-verdigris-dark"
              >
                Get started
                <ArrowRight className="h-4 w-4" />
              </Link>
            </>
          )}
        </nav>
      </header>

      <main className="flex flex-1 flex-col">
        <section className="bg-grid relative overflow-hidden border-b border-line">
          <div className="pointer-events-none absolute inset-0 bg-gradient-to-b from-paper/40 via-paper/80 to-paper" />
          <div className="relative mx-auto grid w-full max-w-6xl items-center gap-12 px-6 py-16 lg:grid-cols-[1.05fr_0.95fr] lg:py-24">
            <div className="flex flex-col gap-6">
              <span className="inline-flex w-fit items-center gap-2 rounded-full border border-line bg-surface px-3 py-1 font-mono text-[11px] uppercase tracking-widest text-graphite">
                <span className="h-1.5 w-1.5 rounded-full bg-verdigris" />
                S3-style object storage
              </span>
              <h1 className="text-balance font-display text-4xl font-semibold leading-[1.05] tracking-tight text-ink sm:text-5xl lg:text-6xl">
                Every file,
                <br />
                catalogued in a bucket.
              </h1>
              <p className="max-w-lg text-base leading-relaxed text-graphite sm:text-lg">
                ObjexRG is object storage with a working web drive. Create
                buckets, push files, mint share links and control who can read
                what — all behind secure, rotating sessions.
              </p>
              <div className="flex flex-wrap items-center gap-3">
                <Link
                  to={cta}
                  className="inline-flex items-center gap-2 rounded-lg bg-verdigris px-6 py-3 text-sm font-medium text-white hover:bg-verdigris-dark"
                >
                  {authenticated ? "Open Drive" : "Create your first bucket"}
                  <ArrowRight className="h-4 w-4" />
                </Link>
                <a
                  href="#features"
                  className="rounded-lg border border-line bg-surface px-6 py-3 text-sm font-medium text-ink hover:border-verdigris hover:text-verdigris"
                >
                  How it works
                </a>
              </div>
              <dl className="mt-2 flex flex-wrap gap-x-8 gap-y-3 border-t border-line pt-6">
                <Stat value="Buckets" label="scoped per account" />
                <Stat value="1 GB" label="per upload" />
                <Stat value="JWT" label="rotating sessions" />
              </dl>
            </div>

            <BucketPreview />
          </div>
        </section>

        <section id="features" className="mx-auto w-full max-w-6xl px-6 py-16 lg:py-24">
          <div className="mb-10 flex flex-col gap-3">
            <span className="font-mono text-[11px] uppercase tracking-widest text-copper">
              Capabilities
            </span>
            <h2 className="text-balance font-display text-2xl font-semibold text-ink sm:text-3xl">
              A drive built around objects, not folders.
            </h2>
          </div>
          <div className="grid grid-cols-1 border-t border-l border-line sm:grid-cols-2 lg:grid-cols-3">
            {features.map((feature) => (
              <article
                key={feature.tag}
                className="flex flex-col gap-3 border-r border-b border-line bg-surface p-6"
              >
                <div className="flex items-center justify-between">
                  <feature.icon className="h-5 w-5 text-verdigris" />
                  <span className="font-mono text-[10px] uppercase tracking-widest text-muted">
                    {feature.tag}
                  </span>
                </div>
                <h3 className="font-display text-base font-semibold text-ink">
                  {feature.title}
                </h3>
                <p className="text-sm leading-relaxed text-graphite">
                  {feature.description}
                </p>
              </article>
            ))}
            <article className="flex flex-col justify-between gap-3 border-r border-b border-line bg-verdigris p-6 text-white">
              <span className="font-mono text-[10px] uppercase tracking-widest text-white/70">
                Start
              </span>
              <div className="flex flex-col gap-3">
                <h3 className="font-display text-base font-semibold">
                  Ready when you are.
                </h3>
                <Link
                  to={cta}
                  className="inline-flex w-fit items-center gap-2 rounded-lg bg-white px-4 py-2 text-sm font-medium text-verdigris-dark hover:bg-white/90"
                >
                  {authenticated ? "Open Drive" : "Get started"}
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </div>
            </article>
          </div>
        </section>
      </main>

      <footer className="border-t border-line px-6 py-6">
        <div className="mx-auto flex w-full max-w-6xl items-center justify-between font-mono text-[11px] uppercase tracking-widest text-muted">
          <span>ObjexRG</span>
          <span>object storage / web drive</span>
        </div>
      </footer>
    </div>
  )
}

function Stat({ value, label }: { value: string; label: string }) {
  return (
    <div className="flex flex-col gap-0.5">
      <dt className="font-display text-lg font-semibold text-ink">{value}</dt>
      <dd className="font-mono text-[10px] uppercase tracking-widest text-muted">
        {label}
      </dd>
    </div>
  )
}

function CapacityBar() {
  const filled = 14
  return (
    <div className="flex gap-[3px]" aria-hidden="true">
      {Array.from({ length: 24 }).map((_, index) => (
        <span
          key={index}
          className={`h-6 flex-1 rounded-[2px] ${
            index < filled ? "bg-verdigris" : "bg-line"
          }`}
        />
      ))}
    </div>
  )
}

function BucketPreview() {
  return (
    <div className="relative">
      <div className="pointer-events-none absolute -inset-3 rounded-3xl bg-verdigris/5" />
      <div className="relative rounded-2xl border border-line bg-surface shadow-2xl">
        <div className="flex items-center justify-between border-b border-line px-5 py-4">
          <div className="flex items-center gap-2 font-mono text-sm text-ink">
            <FolderTree className="h-4 w-4 text-verdigris" />
            photos<span className="text-muted">/</span>
          </div>
          <span className="rounded-full bg-verdigris-tint px-2 py-0.5 font-mono text-[10px] uppercase tracking-widest text-verdigris-dark">
            live
          </span>
        </div>

        <div className="flex flex-col gap-3 px-5 py-4">
          <div className="flex items-center justify-between font-mono text-[10px] uppercase tracking-widest text-muted">
            <span>capacity</span>
            <span>58%</span>
          </div>
          <CapacityBar />
        </div>

        <ul className="flex flex-col border-t border-line">
          {previewRows.map((row) => (
            <li
              key={row.name}
              className="flex items-center gap-3 border-b border-line px-5 py-3 last:border-b-0"
            >
              <span className="flex h-8 w-8 items-center justify-center rounded-md bg-paper text-graphite">
                {row.kind === "png" ? (
                  <ImageIcon className="h-4 w-4" />
                ) : (
                  <FileText className="h-4 w-4" />
                )}
              </span>
              <span className="flex-1 truncate font-mono text-xs text-ink">
                {row.name}
              </span>
              <span className="font-mono text-[10px] uppercase tracking-widest text-muted">
                {row.kind}
              </span>
              <span className="font-mono text-[11px] text-graphite">
                {row.size}
              </span>
              {row.public ? (
                <span className="rounded-full bg-copper-tint px-2 py-0.5 font-mono text-[10px] uppercase tracking-widest text-copper">
                  public
                </span>
              ) : null}
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}
